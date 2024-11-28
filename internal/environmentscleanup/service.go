package environmentscleanup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dbmodel"

	"github.com/hashicorp/go-multierror"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	kcpNamespace = "kcp-system"
	kebConfigMap = "kcp-kyma-environment-broker"
)

//go:generate mockery --name=GardenerClient --output=automock
type GardenerClient interface {
	List(context context.Context, opts v1.ListOptions) (*unstructured.UnstructuredList, error)
	Get(ctx context.Context, name string, options v1.GetOptions, subresources ...string) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, name string, options v1.DeleteOptions, subresources ...string) error
	Update(ctx context.Context, obj *unstructured.Unstructured, options v1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error)
}

//go:generate mockery --name=BrokerClient --output=automock
type BrokerClient interface {
	Deprovision(instance internal.Instance) (string, error)
}

type Service struct {
	gardenerService GardenerClient
	brokerService   BrokerClient
	k8sClient       client.Client
	instanceStorage storage.Instances
	logger          *log.Logger
	MaxShootAge     time.Duration
	LabelSelector   string
	instances       map[string]internal.Instance
	runtimeCRs      map[string]imv1.Runtime
}

type runtime struct {
	ID        string
	ShootName string
}

type providers struct {
	Providers []struct {
		DomainsInclude []string `yaml:"domainsInclude"`
	} `yaml:"providers"`
}

func NewService(gardenerClient GardenerClient, brokerClient BrokerClient, k8sClient client.Client, instanceStorage storage.Instances, logger *log.Logger, maxShootAge time.Duration, labelSelector string) *Service {
	return &Service{
		gardenerService: gardenerClient,
		brokerService:   brokerClient,
		k8sClient:       k8sClient,
		instanceStorage: instanceStorage,
		logger:          logger,
		MaxShootAge:     maxShootAge,
		LabelSelector:   labelSelector,
	}
}

func (s *Service) Run() error {
	environment, err := s.getEnvironment()
	if err != nil {
		return err
	}
	log.Infof("Current environment: %s", environment)
	if environment != "dev.kyma.ondemand.com" {
		return fmt.Errorf("job must run only in the dev environment, current environment: %s", environment)
	}
	return s.PerformCleanup()
}

func (s *Service) getEnvironment() (string, error) {
	configMap := &coreV1.ConfigMap{}
	err := s.k8sClient.Get(context.Background(), client.ObjectKey{
		Namespace: kcpNamespace,
		Name:      kebConfigMap,
	}, configMap)
	if err != nil {
		return "", fmt.Errorf("while getting ConfigMap %s from namespace %s: %w", kebConfigMap, kcpNamespace, err)
	}

	data, exists := configMap.Data["skrDNSProvidersValues.yaml"]
	if !exists {
		return "", fmt.Errorf("while checking skrDNSProvidersValues.yaml key in ConfigMap %s: key is missing", kebConfigMap)
	}

	var providers providers
	err = yaml.Unmarshal([]byte(data), &providers)
	if err != nil {
		return "", fmt.Errorf("while unmarshalling YAML data from skrDNSProvidersValues.yaml: %w", err)
	}

	if len(providers.Providers) != 1 && len(providers.Providers[0].DomainsInclude) != 1 {
		return "", fmt.Errorf("while validating structure of skrDNSProvidersValues.yaml: expected 1 provider with 1 domain, found %d providers", len(providers.Providers))
	}

	return providers.Providers[0].DomainsInclude[0], nil
}

func (s *Service) PerformCleanup() error {
	runtimesToDelete, runtimeCRsToDelete, shootsToDelete, err := s.getStaleRuntimesByShoots(s.LabelSelector)
	if err != nil {
		s.logger.Error(fmt.Errorf("while getting stale shoots to delete: %w", err))
		return err
	}

	err = s.cleanupRuntimes(runtimesToDelete)
	if err != nil {
		s.logger.Error(fmt.Errorf("while cleaning runtimes: %w", err))
		return err
	}

	err = s.cleanUpRuntimeCRs(runtimeCRsToDelete)
	if err != nil {
		s.logger.Error(fmt.Errorf("while cleaning runtime CRs: %w", err))
		return err
	}

	return s.cleanupShoots(shootsToDelete)
}

func (s *Service) cleanupRuntimes(runtimes []runtime) error {
	s.logger.Infof("Runtimes to process: %+v", runtimes)

	if len(runtimes) == 0 {
		return nil
	}

	return s.cleanUp(runtimes)
}

func (s *Service) cleanupShoots(shoots []unstructured.Unstructured) error {
	// do not log all shoots as previously - too much info
	s.logger.Infof("Number of shoots to process: %+v", len(shoots))

	if len(shoots) == 0 {
		return nil
	}

	for _, shoot := range shoots {
		annotations := shoot.GetAnnotations()
		annotations["confirmation.gardener.cloud/deletion"] = "true"
		shoot.SetAnnotations(annotations)
		_, err := s.gardenerService.Update(context.Background(), &shoot, v1.UpdateOptions{})
		if err != nil {
			s.logger.Error(fmt.Errorf("while annotating shoot with removal confirmation: %w", err))
		}

		err = s.gardenerService.Delete(context.Background(), shoot.GetName(), v1.DeleteOptions{})
		if err != nil {
			s.logger.Error(fmt.Errorf("while cleaning runtimes: %w", err))
		}
	}

	return nil
}

func (s *Service) getStaleRuntimesByShoots(labelSelector string) ([]runtime, []runtime, []unstructured.Unstructured, error) {
	opts := v1.ListOptions{
		LabelSelector: labelSelector,
	}
	shootList, err := s.gardenerService.List(context.Background(), opts)
	if err != nil {
		return []runtime{}, []runtime{}, []unstructured.Unstructured{}, fmt.Errorf("while listing Gardener shoots: %w", err)
	}

	s.instances = make(map[string]internal.Instance)
	instanceList, _, _, err := s.instanceStorage.List(dbmodel.InstanceFilter{})
	if err != nil {
		return []runtime{}, []runtime{}, []unstructured.Unstructured{}, fmt.Errorf("while listing instances: %w", err)
	}
	for _, instance := range instanceList {
		s.instances[instance.InstanceDetails.ShootName] = instance
	}

	s.runtimeCRs = make(map[string]imv1.Runtime)
	var runtimeCRList imv1.RuntimeList
	err = s.k8sClient.List(context.Background(), &runtimeCRList, &client.ListOptions{Namespace: kcpNamespace})
	if err != nil {
		return []runtime{}, []runtime{}, []unstructured.Unstructured{}, fmt.Errorf("while listing runtime CRs: %w", err)
	}
	for _, runtimeCR := range runtimeCRList.Items {
		s.runtimeCRs[runtimeCR.Spec.Shoot.Name] = runtimeCR
	}

	var runtimes []runtime
	var runtimeCRs []runtime
	var shoots []unstructured.Unstructured
	for _, shoot := range shootList.Items {
		shootCreationTimestamp := shoot.GetCreationTimestamp()
		shootAge := time.Since(shootCreationTimestamp.Time)

		if shootAge.Hours() < s.MaxShootAge.Hours() {
			log.Infof("Shoot %q is not older than %f hours with age: %f hours", shoot.GetName(), s.MaxShootAge.Hours(), shootAge.Hours())
			continue
		}

		log.Infof("Shoot %q is older than %f hours with age: %f hours", shoot.GetName(), s.MaxShootAge.Hours(), shootAge.Hours())

		instance, ok := s.instances[shoot.GetName()]
		if ok {
			s.logger.Infof("Found an instance %q for shoot %q", instance.InstanceID, shoot.GetName())
			staleRuntime := runtime{
				ID:        instance.RuntimeID,
				ShootName: shoot.GetName(),
			}
			runtimes = append(runtimes, staleRuntime)
			continue
		}

		runtimeCR, ok := s.runtimeCRs[shoot.GetName()]
		if ok {
			s.logger.Infof("Found a runtime CR %q for shoot %q", runtimeCR.Name, shoot.GetName())
			staleRuntimeCR := runtime{
				ID:        runtimeCR.Name,
				ShootName: shoot.GetName(),
			}
			runtimeCRs = append(runtimeCRs, staleRuntimeCR)
			continue
		}

		s.logger.Infof("Instance and runtime CR not found for shoot %q", shoot.GetName())
		shoots = append(shoots, shoot)
	}

	return runtimes, runtimeCRs, shoots, nil
}

func (s *Service) cleanUp(runtimesToDelete []runtime) error {
	kebInstancesToDelete, err := s.getInstancesForRuntimes(runtimesToDelete)
	if err != nil {
		errMsg := fmt.Errorf("while getting instance IDs for Runtimes: %w", err)
		s.logger.Error(errMsg)
		if !dberr.IsNotFound(err) {
			return errMsg
		}
	}

	result := s.cleanUpKEBInstances(kebInstancesToDelete)

	if result != nil {
		result.ErrorFormat = func(i []error) string {
			var s []string
			for _, v := range i {
				s = append(s, v.Error())
			}
			return strings.Join(s, ", ")
		}
	}

	return result.ErrorOrNil()
}

func (s *Service) getInstancesForRuntimes(runtimesToDelete []runtime) ([]internal.Instance, error) {

	var runtimeIDsToDelete []string
	for _, runtime := range runtimesToDelete {
		runtimeIDsToDelete = append(runtimeIDsToDelete, runtime.ID)
	}

	instances, err := s.instanceStorage.FindAllInstancesForRuntimes(runtimeIDsToDelete)
	if err != nil {
		return []internal.Instance{}, err
	}

	return instances, nil
}

func (s *Service) cleanUpKEBInstances(instancesToDelete []internal.Instance) *multierror.Error {
	var result *multierror.Error

	for _, instance := range instancesToDelete {
		s.logger.Infof("Triggering environment deprovisioning for instance ID %q, runtime ID %q and shoot name %q", instance.InstanceID, instance.RuntimeID, instance.InstanceDetails.ShootName)
		currentErr := s.triggerEnvironmentDeprovisioning(instance)
		if currentErr != nil {
			result = multierror.Append(result, currentErr)
		}
	}

	return result
}

func (s *Service) triggerEnvironmentDeprovisioning(instance internal.Instance) error {
	opID, err := s.brokerService.Deprovision(instance)
	if err != nil {
		err = fmt.Errorf("while triggering deprovisioning for instance ID %q: %w", instance.InstanceID, err)
		s.logger.Error(err)
		return err
	}

	log.Infof("Successfully send deprovision request to Kyma Environment Broker, got operation ID %q", opID)
	return nil
}

func (s *Service) cleanUpRuntimeCRs(runtimeCRs []runtime) error {
	s.logger.Infof("RuntimeCRs to process: %+v", runtimeCRs)

	var result *multierror.Error

	for _, runtime := range runtimeCRs {
		s.logger.Infof("Deleting runtime CR for runtimeID ID %q and shoot name %q", runtime.ID, runtime.ShootName)
		err := s.deleteRuntimeCR(runtime)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	if result != nil {
		result.ErrorFormat = func(i []error) string {
			var s []string
			for _, v := range i {
				s = append(s, v.Error())
			}
			return strings.Join(s, ", ")
		}
	}

	return result.ErrorOrNil()
}

func (s *Service) deleteRuntimeCR(runtime runtime) error {
	var runtimeCR = imv1.Runtime{}
	err := s.k8sClient.Get(context.Background(), client.ObjectKey{Name: runtime.ID, Namespace: kcpNamespace}, &runtimeCR)
	if err != nil {
		s.logger.Error(fmt.Errorf("while getting runtime CR for runtime ID %q: %w", runtime.ID, err))
		return nil
	}

	if runtime.ShootName != runtimeCR.Spec.Shoot.Name {
		s.logger.Error(fmt.Errorf("gardener shoot name %q does not match runtime CR shoot name %q", runtime.ShootName, runtimeCR.Spec.Shoot.Name))
		return nil
	}

	err = s.k8sClient.Delete(context.Background(), &runtimeCR)
	if err != nil {
		s.logger.Error(fmt.Errorf("while deleting runtime CR for runtime ID %q: %w", runtime.ID, err))
		return err
	}

	s.logger.Infof("Successfully deleted runtime CR for runtimeID ID %q", runtime.ID)
	return nil
}

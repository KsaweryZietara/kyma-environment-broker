package edp

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kyma-project/kyma-environment-broker/internal/httputil"
	"github.com/stretchr/testify/assert"
)

const (
	subAccountID  = "72b83910-c2dc-415b-b95d-960cc45b36abx"
	environment   = "test"
	testNamespace = "testNs"
)

func TestClient_CreateDataTenant(t *testing.T) {
	// given
	testServer := fixHTTPServer(t)
	defer testServer.Close()

	config := Config{
		AdminURL:  testServer.URL,
		Namespace: testNamespace,
	}
	client := NewClient(config)
	client.setHttpClient(testServer.Client())

	// when
	err := client.CreateDataTenant(DataTenantPayload{
		Name:        subAccountID,
		Environment: environment,
	}, fixLogger())

	// then
	assert.NoError(t, err)

	response, err := testServer.Client().Get(fmt.Sprintf("%s/namespaces/%s/dataTenants/%s/%s", testServer.URL, testNamespace, subAccountID, environment))
	assert.NoError(t, err)

	var dt DataTenantItem
	err = json.NewDecoder(response.Body).Decode(&dt)
	assert.NoError(t, err)
	assert.Equal(t, subAccountID, dt.Name)
	assert.Equal(t, environment, dt.Environment)
	assert.Equal(t, testNamespace, dt.Namespace.Name)
}

func TestClient_DeleteDataTenant(t *testing.T) {
	// given
	testServer := fixHTTPServer(t)
	defer testServer.Close()

	config := Config{
		AdminURL:  testServer.URL,
		Namespace: testNamespace,
	}
	client := NewClient(config)
	client.setHttpClient(testServer.Client())

	err := client.CreateDataTenant(DataTenantPayload{
		Name:        subAccountID,
		Environment: environment,
	}, fixLogger())
	assert.NoError(t, err)

	// when
	err = client.DeleteDataTenant(subAccountID, environment, fixLogger())

	// then
	assert.NoError(t, err)

	response, err := testServer.Client().Get(fmt.Sprintf("%s/namespaces/%s/dataTenants/%s/%s", testServer.URL, testNamespace, subAccountID, environment))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestClient_CreateMetadataTenant(t *testing.T) {
	// given
	testServer := fixHTTPServer(t)
	defer testServer.Close()

	config := Config{
		AdminURL:  testServer.URL,
		Namespace: testNamespace,
	}
	client := NewClient(config)
	client.setHttpClient(testServer.Client())

	// when
	err := client.CreateMetadataTenant(subAccountID, environment, MetadataTenantPayload{Key: "tK", Value: "tV"}, fixLogger())
	assert.NoError(t, err)

	err = client.CreateMetadataTenant(subAccountID, environment, MetadataTenantPayload{Key: "tK2", Value: "tV2"}, fixLogger())
	assert.NoError(t, err)

	// then
	assert.NoError(t, err)

	data, err := client.GetMetadataTenant(subAccountID, environment)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestClient_DeleteMetadataTenant(t *testing.T) {
	// given
	key := "tK"
	testServer := fixHTTPServer(t)
	defer testServer.Close()

	config := Config{
		AdminURL:  testServer.URL,
		Namespace: testNamespace,
	}
	client := NewClient(config)
	client.setHttpClient(testServer.Client())

	err := client.CreateMetadataTenant(subAccountID, environment, MetadataTenantPayload{Key: key, Value: "tV"}, fixLogger())
	assert.NoError(t, err)

	// when
	err = client.DeleteMetadataTenant(subAccountID, environment, key, fixLogger())

	// then
	assert.NoError(t, err)

	data, err := client.GetMetadataTenant(subAccountID, environment)
	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestClient_UpdateMetadata(t *testing.T) {
	// given
	key := "tK"
	testServer := fixHTTPServer(t)
	defer testServer.Close()

	config := Config{
		AdminURL:  testServer.URL,
		Namespace: testNamespace,
	}
	client := NewClient(config)
	client.setHttpClient(testServer.Client())

	err := client.CreateMetadataTenant(subAccountID, environment, MetadataTenantPayload{Key: key, Value: "tV"}, fixLogger())
	assert.NoError(t, err)

	// when
	err = client.UpdateMetadataTenant(subAccountID, environment, key, "newValue", fixLogger())

	assert.NoError(t, err)
	data, err := client.GetMetadataTenant(subAccountID, environment)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "newValue", data[0].Value)
}

func fixHTTPServer(t *testing.T) *httptest.Server {
	r := httputil.NewRouter()
	srv := newServer(t)

	r.HandleFunc("POST /namespaces/{namespace}/dataTenants", srv.createDataTenant)
	r.HandleFunc("DELETE /namespaces/{namespace}/dataTenants/{name}/{env}", srv.deleteDataTenant)

	r.HandleFunc("POST /namespaces/{namespace}/dataTenants/{name}/{env}/metadata", srv.createMetadata)
	r.HandleFunc("PUT /namespaces/{namespace}/dataTenants/{name}/{env}/metadata/{key}", srv.updateMetadata)
	r.HandleFunc("GET /namespaces/{namespace}/dataTenants/{name}/{env}/metadata", srv.getMetadata)
	r.HandleFunc("DELETE /namespaces/{namespace}/dataTenants/{name}/{env}/metadata/{key}", srv.deleteMetadata)

	// enpoints use only for test (exist in real EDP)
	r.HandleFunc("GET /namespaces/{namespace}/dataTenants/{name}/{env}", srv.getDataTenants)

	return httptest.NewServer(r)
}

type server struct {
	t          *testing.T
	metadata   []MetadataItem
	dataTenant map[string][]byte
}

func newServer(t *testing.T) *server {
	return &server{
		t:          t,
		metadata:   make([]MetadataItem, 0),
		dataTenant: make(map[string][]byte, 0),
	}
}

func (s *server) checkNamespace(w http.ResponseWriter, r *http.Request) (string, bool) {
	namespace := r.PathValue("namespace")
	if len(namespace) == 0 {
		s.t.Error("key namespace doesn't exist")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return "", false
	}
	if namespace != testNamespace {
		s.t.Errorf("key namespace is not equal to %s", testNamespace)
		w.WriteHeader(http.StatusNotFound)
		return namespace, false
	}

	return namespace, true
}

func (s *server) fetchNameAndEnv(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	name := r.PathValue("name")
	env := r.PathValue("env")

	if len(name) == 0 || len(env) == 0 {
		s.t.Error("one of the required key doesn't exist")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return "", "", false
	}

	return name, env, true
}

func (s *server) createDataTenant(w http.ResponseWriter, r *http.Request) {
	ns, ok := s.checkNamespace(w, r)
	if !ok {
		return
	}

	var dt DataTenantPayload
	err := json.NewDecoder(r.Body).Decode(&dt)
	if err != nil {
		s.t.Errorf("cannot read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dti := DataTenantItem{
		Namespace: NamespaceItem{
			Name: ns,
		},
		Name:        dt.Name,
		Environment: dt.Environment,
	}

	data, err := json.Marshal(dti)
	if err != nil {
		s.t.Errorf("wrong request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.dataTenant[dt.Name] = data
	w.WriteHeader(http.StatusCreated)
}

func (s *server) deleteDataTenant(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.checkNamespace(w, r); !ok {
		return
	}
	name, _, ok := s.fetchNameAndEnv(w, r)
	if !ok {
		return
	}

	for dtName := range s.dataTenant {
		if dtName == name {
			delete(s.dataTenant, dtName)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// EDP server return 204 if dataTenant not exist already
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) createMetadata(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.checkNamespace(w, r); !ok {
		return
	}
	name, env, ok := s.fetchNameAndEnv(w, r)
	if !ok {
		return
	}

	var item MetadataItem
	err := json.NewDecoder(r.Body).Decode(&item)
	item.DataTenant = DataTenantItem{
		Namespace: NamespaceItem{
			Name: testNamespace,
		},
		Name:        name,
		Environment: env,
	}
	if err != nil {
		s.t.Errorf("cannot decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.metadata = append(s.metadata, item)
	w.WriteHeader(http.StatusCreated)
}

func (s *server) updateMetadata(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.checkNamespace(w, r); !ok {
		return
	}
	key := r.PathValue("key")

	var item MetadataItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		s.t.Errorf("cannot decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, existingItem := range s.metadata {
		if key == existingItem.Key {
			s.metadata[i].Value = item.Value
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) deleteMetadata(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.checkNamespace(w, r); !ok {
		return
	}

	key := r.PathValue("key")
	if len(key) == 0 {
		s.t.Error("key doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newData := make([]MetadataItem, 0)
	for _, item := range s.metadata {
		if item.Key == key {
			continue
		}
		newData = append(newData, item)
	}

	s.metadata = newData
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getMetadata(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(s.metadata)
	if err != nil {
		s.t.Errorf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) getDataTenants(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.checkNamespace(w, r); !ok {
		return
	}
	name, _, ok := s.fetchNameAndEnv(w, r)
	if !ok {
		s.t.Error("cannot find name/env query parameters")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, ok := s.dataTenant[name]; !ok {
		s.t.Logf("dataTenant with name %s not exist", name)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := w.Write(s.dataTenant[name])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// setHttpClient auxiliary method of testing to get rid of oAuth client wrapper
func (c *Client) setHttpClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func fixLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

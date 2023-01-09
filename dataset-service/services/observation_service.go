package services

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	common_config "github.com/caraml-dev/timber/common/config"
	"github.com/caraml-dev/timber/common/log"
	timberv1 "github.com/caraml-dev/timber/dataset-service/api"
	"github.com/caraml-dev/timber/dataset-service/config"
	"github.com/caraml-dev/timber/dataset-service/models"
	osConfig "github.com/caraml-dev/timber/observation-service/config"
)

const (
	// HELM_DRIVER is the storage backend to be used. Documentation details: https://helm.sh/docs/topics/advanced/#storage-backends
	HELM_DRIVER = "secret"
	// RELEASE_NAME is the helm release name prefix to be used when deploying Observation Service.
	RELEASE_NAME = "observation-service"
)

// ObservationService provides a set of methods to interact with the MLP APIs
type ObservationService interface {
	// CreateService creates new Observation Service Helm release and returns ID of created Observation Service
	CreateService(projectName string, config *timberv1.ObservationServiceConfig) (*timberv1.ObservationServiceResponse, error)
	// UpdateService updates existing Observation Service Helm release and returns ID of updated Observation Service
	UpdateService(
		projectName string,
		observationServiceID int,
		config *timberv1.ObservationServiceConfig,
	) (*timberv1.ObservationServiceResponse, error)
}

type observationService struct {
	services *Services

	gcpProject               string
	deploymentConfig         common_config.DeploymentConfig
	observationServiceConfig config.ObservationServiceConfig
}

// NewObservationService instantiates ObservationService
func NewObservationService(
	services *Services,
	deploymentConfig common_config.DeploymentConfig,
	observationServiceConfig config.ObservationServiceConfig,
) ObservationService {
	return &observationService{
		services:                 services,
		gcpProject:               observationServiceConfig.GCPProject,
		deploymentConfig:         deploymentConfig,
		observationServiceConfig: observationServiceConfig,
	}
}

func (o *observationService) CreateService(
	caramlProjectName string,
	config *timberv1.ObservationServiceConfig,
) (*timberv1.ObservationServiceResponse, error) {
	// Read chart
	chart, err := readChart()
	if err != nil {
		return nil, err
	}

	// Check if chart is installable
	validInstallableChart, err := isChartInstallable(chart)
	if !validInstallableChart {
		log.Info(err)
	}

	// Retrieve computed chart values based on default values and request body values
	updatedChartValues, err := getChartValues(config, o.gcpProject, caramlProjectName, o.deploymentConfig, o.observationServiceConfig)
	if err != nil {
		return nil, err
	}

	// Generate configuration required to run helm installation, this is dependent on storage backend.
	settings := cli.New()
	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), caramlProjectName, HELM_DRIVER, log.Infof)
	if err != nil {
		return nil, err
	}

	// Postfix helm release name with provided Service name as a CaraML project could have multiple different services
	projectReleaseName := fmt.Sprintf("%s-%s", RELEASE_NAME, config.GetServiceName())

	// Initialize helm installation
	installation := action.NewInstall(actionConfig)
	installation.ReleaseName = projectReleaseName
	installation.Namespace = caramlProjectName
	installation.CreateNamespace = true

	// Trigger helm installation
	release, err := installation.Run(chart, updatedChartValues)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("%s helm release installation (version %d) underway! Status: %s", release.Name, release.Version, release.Info.Status)
	// Pretty print Helm release YAML manifest with fmt.Println
	fmt.Println(release.Manifest)

	// TODO: Retrieve Observation Service ID from DB
	resp := &timberv1.ObservationServiceResponse{
		Id: uuid.New().String(),
	}

	return resp, nil
}

func (o *observationService) UpdateService(
	projectName string,
	observationServiceID int,
	config *timberv1.ObservationServiceConfig,
) (*timberv1.ObservationServiceResponse, error) {
	// Read chart
	chart, err := readChart()
	if err != nil {
		return nil, err
	}

	// Check if chart is installable
	validInstallableChart, err := isChartInstallable(chart)
	if !validInstallableChart {
		log.Info(err)
	}

	// Retrieve computed chart values based on default values and request body values
	updatedChartValues, err := getChartValues(config, o.gcpProject, projectName, o.deploymentConfig, o.observationServiceConfig)
	if err != nil {
		return nil, err
	}

	// Generate configuration required to run helm upgrade
	settings := cli.New()
	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), projectName, HELM_DRIVER, log.Infof)
	if err != nil {
		return nil, err
	}

	// Initialize helm upgrade
	upgrade := action.NewUpgrade(actionConfig)
	upgrade.Namespace = projectName

	// Trigger helm upgrade
	// TODO: Get project release name based on provided Observation Service ID
	projectReleaseName := fmt.Sprintf("%s-%s", RELEASE_NAME, config.GetServiceName())
	release, err := upgrade.Run(projectReleaseName, chart, updatedChartValues)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("%s helm release upgrade (version %d) underway! Status: %s", release.Name, release.Version, release.Info.Status)
	// Pretty print Helm release YAML manifest with fmt.Println
	fmt.Println(release.Manifest)

	// TODO: Retrieve Observation Service ID from DB
	resp := &timberv1.ObservationServiceResponse{
		Id: strconv.Itoa(observationServiceID),
	}

	return resp, nil
}

func readChart() (*chart.Chart, error) {
	// TODO: Publish helm chart and reference published chart
	// Locate directory of Observation Service helm chart
	_, filename, _, _ := runtime.Caller(0)
	repositoryRoot := filepath.Dir(filepath.Dir(path.Dir(filename)))
	chartDir := fmt.Sprintf("%s/infra/charts/observation-service", repositoryRoot)
	chart, err := loader.LoadDir(chartDir)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func getChartValues(
	config *timberv1.ObservationServiceConfig,
	gcpProject string,
	caramlProjectName string,
	deploymentConfig common_config.DeploymentConfig,
	observationServiceConfig config.ObservationServiceConfig,
) (map[string]interface{}, error) {
	// Initialize and set default values
	values := &models.ObservationServiceHelmValues{}
	values = setDefaultValues(values, gcpProject, caramlProjectName, config.GetServiceName(), deploymentConfig, observationServiceConfig)

	// Compute LogConsumerConfig
	values, err := setLogConsumerConfig(values, config)
	if err != nil {
		return nil, err
	}

	// Compute LogProducerConfig
	values, err = setLogProducerConfig(values, config, gcpProject, caramlProjectName)
	if err != nil {
		return nil, err
	}

	// Compute Fluentd
	if config.GetSink().GetType() == timberv1.ObservationServiceDataSinkType_OBSERVATION_SERVICE_DATA_SINK_TYPE_FLUENTD {
		values, err = setFluentdConfig(values, config)
		if err != nil {
			return nil, err
		}
	}

	// Convert type of values for merging
	var interfaceValues map[string]interface{}
	byteArr, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byteArr, &interfaceValues)
	if err != nil {
		return nil, err
	}

	return interfaceValues, nil
}

// setDefaultValues configures default values handled by Dataset Service for Observation Service deployments
func setDefaultValues(
	values *models.ObservationServiceHelmValues,
	gcpProject string,
	caramlProjectName string,
	serviceName string,
	deploymentConfig common_config.DeploymentConfig,
	observationServiceConfig config.ObservationServiceConfig,
) *models.ObservationServiceHelmValues {
	// --- Observation Service configs --- //
	// Image
	values.ObservationServiceConfig.Image.Tag = observationServiceConfig.ObservationServiceImageTag
	// ApiConfig
	values.ObservationServiceConfig.ApiConfig.Port = 9001
	values.ObservationServiceConfig.ApiConfig.DeploymentConfig = deploymentConfig
	values.ObservationServiceConfig.ApiConfig.DeploymentConfig.ProjectName = caramlProjectName
	values.ObservationServiceConfig.ApiConfig.DeploymentConfig.ServiceName = fmt.Sprintf("%s-%s", RELEASE_NAME, serviceName)
	// Environment Variables
	envVars := []models.Env{
		{
			Name:  "GOOGLE_APPLICATION_CREDENTIALS",
			Value: "/etc/gcp_service_account/service-account.json",
		},
	}
	values.ObservationServiceConfig.ExtraEnvs = envVars
	// Resources
	values.ObservationServiceConfig.Resources.Requests.CPU = "1"
	values.ObservationServiceConfig.Resources.Requests.Memory = "512Mi"
	values.ObservationServiceConfig.Resources.Limits.CPU = "1"
	values.ObservationServiceConfig.Resources.Limits.Memory = "1Gi"
	// Autoscaling
	values.ObservationServiceConfig.Autoscaling.Enabled = false

	// --- Fluentd configs --- //
	// Image
	values.FluentdConfig.Image.Tag = observationServiceConfig.FluentdImageTag
	// Environment Variables
	envVars = []models.Env{
		{
			Name:  "FLUENTD_WORKER_COUNT",
			Value: "1",
		},
		{
			Name:  "FLUENTD_LOG_LEVEL",
			Value: "debug",
		},
		{
			Name:  "FLUENTD_BUFFER_LIMIT",
			Value: "1g",
		},
		{
			Name:  "FLUENTD_FLUSH_INTERVAL_SECONDS",
			Value: "30",
		},
		{
			Name:  "FLUENTD_LOG_PATH",
			Value: "/cache/log",
		},
		{
			Name:  "FLUENTD_TAG",
			Value: "observation-service",
		},
		{
			Name:  "FLUENTD_GCP_JSON_KEY_PATH",
			Value: "/etc/gcp_service_account/service-account.json",
		},
		{
			Name:  "FLUENTD_GCP_PROJECT",
			Value: gcpProject,
		},
		{
			Name:  "FLUENTD_BQ_DATASET",
			Value: caramlProjectName,
		},
		{
			Name:  "FLUENTD_BQ_TABLE",
			Value: fmt.Sprintf("%s_observation_log", caramlProjectName),
		},
	}
	values.FluentdConfig.ExtraEnvs = envVars
	// Resources
	values.FluentdConfig.Resources.Requests.CPU = "1"
	values.FluentdConfig.Resources.Requests.Memory = "512Mi"
	values.FluentdConfig.Resources.Limits.CPU = "1"
	values.FluentdConfig.Resources.Limits.Memory = "1Gi"
	// Autoscaling
	values.FluentdConfig.Autoscaling.Enabled = false

	return values
}

// setFluentdConfig configures custom values for Fluentd Deployment
func setFluentdConfig(
	values *models.ObservationServiceHelmValues,
	config *timberv1.ObservationServiceConfig,
) (*models.ObservationServiceHelmValues, error) {
	return values, nil
}

// setLogConsumerConfig configures custom values for LogConsumerConfig
func setLogConsumerConfig(
	values *models.ObservationServiceHelmValues,
	config *timberv1.ObservationServiceConfig,
) (*models.ObservationServiceHelmValues, error) {
	switch config.GetSource().GetType() {
	case timberv1.ObservationServiceDataSourceType_OBSERVATION_SERVICE_DATA_SOURCE_TYPE_EAGER:
		return nil, fmt.Errorf("source type (eager) is currently unsupported")
	case timberv1.ObservationServiceDataSourceType_OBSERVATION_SERVICE_DATA_SOURCE_TYPE_KAFKA:
		kafkaConfig := config.GetSource().GetKafkaConfig()
		values.ObservationServiceConfig.ApiConfig.LogConsumerConfig.Kind = "kafka"
		values.ObservationServiceConfig.ApiConfig.LogConsumerConfig.KafkaConfig = models.GetKafkaConfigModel(kafkaConfig)
	case timberv1.ObservationServiceDataSourceType_OBSERVATION_SERVICE_DATA_SOURCE_TYPE_UNSPECIFIED:
		log.Infof("No source type specified for Observation Service deployment")
	default:
		return nil, fmt.Errorf("invalid source type (%s) was provided", config.GetSource().GetType())
	}

	return values, nil
}

// setLogProducerConfig configures custom values for LogProducerConfig
func setLogProducerConfig(
	values *models.ObservationServiceHelmValues,
	config *timberv1.ObservationServiceConfig,
	gcpProject string,
	projectName string,
) (*models.ObservationServiceHelmValues, error) {
	switch config.GetSink().GetType() {
	case timberv1.ObservationServiceDataSinkType_OBSERVATION_SERVICE_DATA_SINK_TYPE_FLUENTD:
		fluentdConfig := config.GetSink().GetFluentdConfig()
		values.ObservationServiceConfig.ApiConfig.LogProducerConfig.Kind = "fluentd"
		values.ObservationServiceConfig.ApiConfig.LogProducerConfig.FluentdConfig = models.GetFluentdConfigModel(fluentdConfig, projectName)
		values.ObservationServiceConfig.ApiConfig.LogProducerConfig.FluentdConfig.BQConfig = &osConfig.BQConfig{
			Project: gcpProject,
			Dataset: projectName,
			Table:   fmt.Sprintf("%s_observation_log", projectName),
		}
	case timberv1.ObservationServiceDataSinkType_OBSERVATION_SERVICE_DATA_SINK_TYPE_KAFKA:
		kafkaConfig := config.GetSink().GetKafkaConfig()
		values.ObservationServiceConfig.ApiConfig.LogProducerConfig.Kind = "kafka"
		values.ObservationServiceConfig.ApiConfig.LogProducerConfig.KafkaConfig = models.GetKafkaConfigModel(kafkaConfig)
	case timberv1.ObservationServiceDataSinkType_OBSERVATION_SERVICE_DATA_SINK_TYPE_STDOUT:
		log.Infof("Standard output sink type specified for Observation Service deployment")
	case timberv1.ObservationServiceDataSinkType_OBSERVATION_SERVICE_DATA_SINK_TYPE_NOOP:
		log.Infof("No-Op sink type specified for Observation Service deployment")
	case timberv1.ObservationServiceDataSinkType_OBSERVATION_SERVICE_DATA_SINK_TYPE_UNSPECIFIED:
		log.Infof("No sink type specified for Observation Service deployment")
	default:
		return nil, fmt.Errorf("invalid sink type (%s) was provided", config.GetSink().GetType())
	}

	return values, nil
}

package common

const (
	//MasterNodeConfigurationFiles file name
	MasterNodeConfigurationFiles = "1.1_master_node_configuration_files.yml"
	//APIServer file name
	APIServer = "1.2_api_server.yml"
	//ControllerManager file name
	ControllerManager = "1.3_controller_manager.yml"
	//Scheduler file name
	Scheduler = "1.4_scheduler.yml"
	//Etcd file name
	Etcd = "2.0_etcd.yml"
	//WorkerNodes file name
	WorkerNodes = "4.0_worker_nodes.yml"
	//ControlPlaneConfiguration file name
	ControlPlaneConfiguration = "3.0_control_plane_configuration.yml"
	//GkeWorkerNodes worker spec
	GkeWorkerNodes = "3.0_worker_nodes.yml"
	//GkePolicies policies spec
	GkePolicies = "4.0_policies.yml"
	//GkeManagedServices manages service spec
	GkeManagedServices = "5.0_managed_services.yml"
	//GkeControlPlaneConfiguration file name
	GkeControlPlaneConfiguration = "2.0_control_plane_configuration.yml"
	//Policies file name
	Policies = "5.0_policies.yml"
	//GrepRegex for tests
	GrepRegex = "[^\"]\\S*'"
	//MultiValue for tests
	MultiValue = "MultiValue"
	//SingleValue for tests
	SingleValue = "SingleValue"
	//EmptyValue for test
	EmptyValue = "EmptyValue"
	//NotValidNumber value
	NotValidNumber = "10000"
	//Report arg
	Report = "r"
	//ReportFull arg
	ReportFull = "report"
	//Synopsis help
	Synopsis = "synopsis"
	//BeaconCli Name
	BeaconCli = "kube-beacon"
	//BeaconVersion version
	BeaconVersion = "0.1"
	//IncludeParam param
	IncludeParam = "i="
	//ExcludeParam param
	ExcludeParam = "e="
	//NodeParam param
	NodeParam = "n="
	//BeaconHomeEnvVar Beacon Home env var
	BeaconHomeEnvVar = "BEACON_HOME"
	//KubeBeacon binary name
	KubeBeacon = "kube-beacon"
	//RootUser process user owner
	RootUser = "root"
	//NonApplicableTest test is not applicable
	NonApplicableTest = "non_applicable"
	//ManualTest test can only be manual executed
	ManualTest = "manual"
	//K8sBenchAuditResultHook hook name
	K8sBenchAuditResultHook = "K8sBenchAuditResultHook"
)

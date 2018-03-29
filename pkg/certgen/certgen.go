package certgen

var DeploymentName string

// SetDeploymentName saves the name of the deployment
func SetDeploymentName(deploymentname string) {
          DeploymentName = deploymentname
}

// GetDeploymentName - Returns the name of the current Deployment
func GetDeploymentName() string {
     return DeploymentName
}




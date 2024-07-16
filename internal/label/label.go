package label

const (
	KeyComponent = "app.kubernetes.io/component"
	KeyCreatedBy = "app.kubernetes.io/created-by"
	KeyInstance  = "app.kubernetes.io/instance"
	KeyManagedBy = "app.kubernetes.io/managed-by"
	KeyName      = "app.kubernetes.io/name"
	KeyPartOf    = "app.kubernetes.io/part-of"
	KeyDashboard = "kyma-project.io/dashboard"

	ValueCompanionBackend = "kyma-companion-backend"
	ValueCompanion        = "companion"
	ValueControllerName   = "kyma-companion-manager"
)

func GetCommonLabels(name string) map[string]string {
	return map[string]string{
		KeyInstance:  name,
		KeyName:      name,
		KeyDashboard: ValueCompanion,
		KeyComponent: ValueCompanion,
		KeyCreatedBy: ValueControllerName,
		KeyManagedBy: ValueControllerName,
		KeyPartOf:    name,
	}
}

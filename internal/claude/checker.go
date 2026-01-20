// Package claude proporciona un wrapper para usar Claude CLI directamente.
//
// Este paquete elimina la dependencia de APIs externas y usa el CLI de Claude
// instalado localmente para generar configuraciones de proyectos.
package claude

// CheckInstalled verifica si claude CLI está instalado y es accesible.
//
// Ejecuta "claude -v" y verifica que la salida sea válida.
// Retorna error si el CLI no está instalado o no es accesible.
func CheckInstalled() error {
	wrapper := NewCLIWrapper()
	return wrapper.CheckInstalled()
}

// GetVersion retorna la versión de Claude CLI instalada.
//
// Retorna la versión como string, o error si no se puede obtener.
func GetVersion() (string, error) {
	wrapper := NewCLIWrapper()
	return wrapper.GetVersion()
}

// Deprecated: Use NewCLIWrapper().SendMessage() directamente.
// Esta función se mantiene por compatibilidad pero se considera deprecada.
func CheckCLI() error {
	return CheckInstalled()
}

# Skill: Desarrollador de Casos de Uso (usecase-developer)

## Propósito
Implementar la lógica de negocio pura y orquestar el flujo de datos dentro de la aplicación.

## Responsabilidades
- Crear clases de caso de uso en `src/application/usecases/` (un archivo por operación/acción).
- Inyectar dependencias mediante interfaces (inversión de control) a través del constructor.
- Implementar la lógica de negocio central, validando reglas de dominio.
- Orquestar llamadas a repositorios, servicios externos y el bus de eventos.
- Retornar siempre DTOs de salida, nunca entidades de infraestructura o base de datos.
- Mantener los casos de uso pequeños, enfocados y con una única responsabilidad.

## Skills Tecnológicas Clave
- **typescript**: Para la implementación de la lógica con tipado fuerte.
- **node-js**: Para el manejo eficiente de la asincronía y el flujo de ejecución.

# Skill: Especialista en Infraestructura (infra-specialist)

## Propósito
Conectar el dominio con el mundo exterior (Base de Datos, APIs externas, Rutas, Controladores) implementando adaptadores y controladores.

## Responsabilidades
- Implementar repositorios de TypeORM en `src/infrastructure/db/repositories/` que cumplan con las interfaces del dominio.
- Crear controladores en `src/infrastructure/controllers/` que orquesten los casos de uso.
- Definir y configurar rutas en `src/infrastructure/routes/` aplicando middlewares necesarios (`isAuth`, permisos, etc.).
- Validar inputs de entrada mediante **Zod** y transformar datos con utilidades como `toXParams`.
- Gestionar la respuesta HTTP utilizando `HttpResponseHandler.ok|error`.
- Configurar adaptadores para servicios externos (Storage, Email, Colas, etc.).

## Skills Tecnológicas Clave
- **express-js**: Para la gestión de rutas y middlewares.
- **api-rest**: Para el diseño semántico de los endpoints.
- **zod**: Para la validación estricta de esquemas de entrada.
- **typescript**: Para el tipado de controladores y adaptadores.

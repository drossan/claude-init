# Skill: Arquitecto de Dominio (domain-expert)

## Propósito
Definir Entidades, DTOs e Interfaces de repositorio siguiendo la Arquitectura Hexagonal y los principios de diseño orientado al dominio.

## Responsabilidades
- Crear entidades en `src/domain/entities/` con validación interna robusta.
- Definir DTOs estrictos en `src/domain/dto/` utilizando `Readonly` y tipos precisos.
- Diseñar contratos (interfaces) de repositorio en `src/domain/interfaces/`.
- Asegurar la total independencia de la capa de dominio respecto a detalles de infraestructura o frameworks.
- Implementar métodos de formato y transformación en las entidades (e.g., `.toDTO()`, `.formatDTO()`).

## Skills Tecnológicas Clave
- **typescript**: Para el tipado estricto de entidades, DTOs e interfaces.

---
name: plan-manage
description: Gestiona los planes de desarrollo del proyecto claude-init CLI. Permite ver el estado, retomar trabajos pausados, y archivar planes completados.
usage: "plan-manage <status|resume|finish> [plan-name]"
---

# Comando: Gestión de Planes (Plan Manage)

Este comando gestiona los planes de desarrollo del proyecto.

## Subcomandos

### status
Muestra el estado actual de todos los planes activos.

```bash
claude-init plan-manage status
```

Salida ejemplo:
```
Active Plans:
  - implement-ai-client [IN_PROGRESS]
    - Implementar Claude API client: DONE
    - Implementar OpenAI API client: IN_PROGRESS
    - Implementar z.ai API client: PENDING
```

### resume
Retoma un plan pausado o archivado.

```bash
claude-init plan-manage resume <plan-name>
```

- Mueve el plan de `completed/` a `active/`
- Restaura la sesión activa
- Muestra el progreso actual

### finish
Marca un plan como completado y lo archiva.

```bash
claude-init plan-manage finish <plan-name>
```

- Verifica que todas las tareas estén completas
- Mueve el plan a `completed/`
- Archiva la sesión
- Añade metadata (fecha de finalización, resumen)

## Estructura de Directorios

```
.claude/
├── plans/
│   ├── active/
│   │   └── implement-ai-client.md
│   └── completed/
│       └── initial-setup.md
└── sessions/
    ├── active/
    │   └── implement-ai-client.md
    └── completed/
        └── initial-setup.md
```

## Formato de Plan Completado

```markdown
# Implementar AI Client

## Completado: 2024-01-15

## Resumen
Se implementaron los clientes de API para Claude, OpenAI y z.ai,
permitiendo que el CLI se conecte a diferentes proveedores de IA.

## Cambios Principales
- [x] Implementar Claude API client
- [x] Implementar OpenAI API client
- [x] Implementar z.ai API client
- [x] Añadir tests unitarios
- [x] Añadir documentación

## Commit
abc123def - Implement AI clients for Claude, OpenAI, and z.ai
```

## Reglas Críticas

- **Verificación**: Verificar que todas las tareas estén completas antes de archivar
- **Metadata**: Añadir fecha y resumen al archivar
- **Limpieza**: Limpiar sesiones activas al finalizar

---

¿Qué acción deseas realizar?

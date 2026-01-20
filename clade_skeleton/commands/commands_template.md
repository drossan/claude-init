---
name: { command-name }
version: 1.0.0
author: { team/person }
description: { Descripción breve, clara y orientada a resultado }
usage: "{command-name} [args] [optional-context]"
type: { planning | executable | meta }
writes_code: false
creates_plan: { true | false }
requires_approval: { true | false }
dependencies: [ other-commands, mcps ]
---

# Comando: {Command Title}

## Objetivo

{Descripción detallada del propósito del command, su alcance y el resultado esperado.
Debe dejar claro si:

- Genera un plan
- Ejecuta validaciones
- Orquesta otros commands
- Requiere aprobación manual}

## Contexto Requerido del Usuario

Lista explícita de información necesaria antes de ejecutar:

- [ ] Descripción de la feature/problema
- [ ] Alcance temporal (sprint, versión, milestone)
- [ ] Restricciones técnicas (stack, librerías prohibidas, compatibilidad)
- [ ] Criterios de aceptación
- [ ] Nivel de prioridad (crítico, alta, media, baja)

## Análisis Inicial (Obligatorio)

Antes de cualquier acción, el command debe evaluar:

- Alcance de la tarea
- Impacto en el código
- Riesgos técnicos
- Dependencias (internas y externas)
- Necesidad de planificación
- Commands o MCPs auxiliares necesarios
- **Agentes y skills óptimos para la tarea**
- Conflictos con planes pendientes

### Pre-ejecución: Checklist Obligatorio

El command debe verificar:

- [ ] ¿La tarea requiere modificar código? → Si sí, generar plan
- [ ] ¿Existen dependencias circulares entre agentes? → Abortar
- [ ] ¿Los skills requeridos están disponibles? → Fallar temprano
- [ ] ¿El contexto del usuario es suficiente? → Solicitar aclaraciones
- [ ] ¿Hay conflictos con planes pendientes? → Avisar y resolver

**Output esperado**: JSON de validación antes de continuar.

```json
{
  "validation_passed": true,
  "risks": [
    "Modifica API pública",
    "Requiere migración de BD"
  ],
  "required_approvals": [
    "tech-lead",
    "security-team"
  ],
  "estimated_complexity": "high",
  "blocking_issues": []
}
```

## Selección de Agentes y Skills (Framework RACI)

El command debe **elegir explícitamente** los agentes y skills más adecuados utilizando el modelo RACI:

- **R** (Responsible): Agente que ejecuta la tarea
- **A** (Accountable): Agente que valida y aprueba
- **C** (Consulted): Skills/MCPs necesarios como soporte
- **I** (Informed): Commands que deben ser notificados

### Criterios de Selección

| Criterio                          | Peso  | Agentes Candidatos                          |
|-----------------------------------|-------|---------------------------------------------|
| Complejidad técnica               | Alta  | `senior-developer`, `architect`             |
| Impacto en arquitectura           | Alta  | `architect`, `tech-lead`                    |
| Tareas repetitivas/automatizables | Media | `junior-developer`, `automation-specialist` |
| Validaciones críticas             | Alta  | `qa-engineer`, `security-expert`            |
| Documentación técnica             | Media | `tech-writer`, `developer`                  |
| Optimización de rendimiento       | Alta  | `performance-engineer`, `senior-developer`  |

> Esto es un ejemplo, en la práctica tendrá que seleccionar los agentes y skills adecuados según el contexto y las
> necesidades específicas del proyecto.

### Ejemplo de Asignación RACI

```yaml
fase_1_diseño:
  responsible: architect
  accountable: tech-lead
  consulted: [ api-design, security-analysis ]
  informed: [ product-manager ]

fase_2_implementacion:
  responsible: backend-developer
  accountable: architect
  consulted: [ code-generation, testing, database-design ]
  informed: [ qa-engineer ]
```

> ⚠️ **La omisión de esta sección invalida el command.**

## Flujo de Trabajo Orquestado

Cada fase debe estar asignada a **un agente concreto** con **responsabilidades claras**.

### 1. {Nombre de la Fase} ({Agente Responsible} | Validado por {Agente Accountable})

**Objetivo**: {Resultado esperado de esta fase}

**Tareas**:

- {Paso concreto y verificable}
- {Paso concreto y verificable}
- {Paso concreto y verificable}

**Asignación**:

- **Agente**: {agent-name}
- **Skills**: `{skill-1}`, `{skill-2}`, `{skill-3}`
- **MCPs**: `{mcp-1}` (opcional)
- **Validador**: {agent-name-validator}

**Criterios de Salida**:

- [ ] {Condición verificable 1}
- [ ] {Condición verificable 2}

---

### 2. {Nombre de la Fase} ({Agente Responsible} | Validado por {Agente Accountable})

**Objetivo**: {Resultado esperado de esta fase}

**Tareas**:

- {Paso concreto y verificable}
- {Paso concreto y verificable}

**Asignación**:

- **Agente**: {agent-name}
- **Skills**: `{skill-1}`, `{skill-2}`
- **Dependencias**: Fase 1 completada
- **Validador**: {agent-name-validator}

**Criterios de Salida**:

- [ ] {Condición verificable 1}
- [ ] {Condición verificable 2}

---

### [N]. {Nombre de la Fase} ({Agente Responsible} | Validado por {Agente Accountable})

**Objetivo**: {Resultado esperado de esta fase}

**Tareas**:

- {Paso concreto y verificable}
- {Paso concreto y verificable}

**Asignación**:

- **Agente**: {agent-name}
- **Skills**: `{skill-1}`, `{skill-2}`
- **Validador**: {agent-name-validator}

**Criterios de Salida**:

- [ ] {Condición verificable 1}
- [ ] {Condición verificable 2}

## Uso de otros Commands y MCPs

{Indicar explícitamente si el command:

- Invoca otros commands (listar cuáles y por qué)
- Utiliza MCPs del proyecto (especificar configuración necesaria)
- Comparte o consume contexto (formato y ubicación)
- Genera eventos para otros commands}

**Ejemplo**:

```yaml
commands_invocados:
  - name: code-analyzer
    trigger: pre-ejecución
    output_required: metrics.json

mcps_utilizados:
  - name: database-schema-validator
    config: .claude/mcp-configs/db-validator.json

contexto_compartido:
  location: .claude/context/shared-state.json
  format: JSON
  consumers: [ qa-automation, deployment-manager ]
```

## Output y Artefactos

| Artefacto                | Ubicación                                     | Formato    | Validador          | Obligatorio     |
|--------------------------|-----------------------------------------------|------------|--------------------|-----------------|
| Plan técnico             | `.claude/plans/{timestamp}-{command-name}.md` | Markdown   | `plan-validator`   | Sí (planning)   |
| Diagrama de arquitectura | `.claude/diagrams/{id}.mmd`                   | Mermaid    | -                  | No              |
| Checklist de validación  | `.claude/checklists/{id}.json`                | JSON       | `schema-validator` | Sí              |
| Reporte de análisis      | `.claude/reports/{id}.md`                     | Markdown   | -                  | Sí (executable) |
| Log de ejecución         | `.claude/logs/{command-name}-{date}.log`      | Plain text | -                  | Sí              |

## Rollback y Cancelación

Si el command falla o el usuario cancela durante la ejecución:

### Procedimiento de Rollback

1. **Detener agentes en curso**: Enviar señal de cancelación a todos los agentes activos
2. **Eliminar artefactos parciales**:
    - Borrar planes incompletos en `.claude/plans/`
    - Limpiar archivos temporales en `.claude/temp/`
3. **Restaurar estado previo**: Si se modificó contexto compartido, revertir a snapshot anterior
4. **Registrar cancelación**:
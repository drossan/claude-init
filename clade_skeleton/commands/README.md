# Guía Técnica para la Creación de Commands en Claude Code

## 1. Propósito de los Commands

Los commands son **orquestadores de tareas**.

Su responsabilidad es analizar una solicitud, definir la estrategia de ejecución, seleccionar los agentes y skills más adecuados, y coordinar el trabajo, **sin modificar código directamente**.

**Un command no ejecuta trabajo técnico, sino que decide quién y cómo debe ejecutarse.**

---

## 2. Tipos de Commands

### 2.1 Planning Commands (con impacto en código)

**Propósito**: Orquestar cambios en el código base.

**Output**: Plan estructurado en `.claude/plans/`

**Ejemplos**: `feature-planner`, `refactor-analyzer`, `migration-coordinator`

#### Reglas obligatorias

❌ **No pueden**:
- Crear código
- Modificar archivos
- Eliminar archivos

✅ **Deben**:
- Analizar la tarea solicitada
- Seleccionar el/los agente(s) y skill(s) que mejor se adapten a la tarea
    - Puede ser uno o varios agentes
    - Cada agente puede tener uno o varios skills
- Definir una estrategia de trabajo clara
- Generar un plan de trabajo en Markdown
- Guardarlo en `.claude/plans/`
- Dejarlo pendiente de aprobación

📌 **Excepción única**: Solo el command `plan-manager` puede iniciar, aprobar o ejecutar planes.

---

### 2.2 Executable Commands (sin impacto en código)

**Propósito**: Validaciones, análisis, reportes y auditorías.

**Output**: Resultados inmediatos (logs, reportes, métricas).

**Ejemplos**: `lint-check`, `security-audit`, `dependency-analyzer`, `performance-profiler`

#### Reglas

✅ Se ejecutan directamente
❌ No generan plan
❌ No pueden modificar archivos

⚠️ **Si durante el análisis se detecta impacto en código, deben abortar y generar un plan en su lugar.**

---

### 2.3 Meta Commands (orquestación de alto nivel)

**Propósito**: Coordinar múltiples commands, planes o flujos de trabajo complejos.

**Output**: Flujo de trabajo compuesto, secuencias de ejecución.

**Ejemplos**: `plan-manager`, `workflow-orchestrator`, `release-coordinator`

#### Restricciones especiales

- **Solo uno por proyecto** (evitar recursión infinita)
- Requieren permisos especiales de ejecución
- Deben implementar detección de ciclos
- Mantienen estado global del proyecto

---

## 3. Template Oficial de Commands

Este template define el contrato técnico obligatorio para todo command.

```markdown
---
name: {command-name}
version: 1.0.0
author: {team/person}
description: {Descripción breve, clara y orientada a resultado}
usage: "{command-name} [args] [optional-context]"
type: {planning | executable | meta}
writes_code: false
creates_plan: {true | false}
requires_approval: {true | false}
dependencies: [other-commands, mcps]
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
  "risks": ["Modifica API pública", "Requiere migración de BD"],
  "required_approvals": ["tech-lead", "security-team"],
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

| Criterio | Peso | Agentes Candidatos |
|----------|------|-------------------|
| Complejidad técnica | Alta | `senior-developer`, `architect` |
| Impacto en arquitectura | Alta | `architect`, `tech-lead` |
| Tareas repetitivas/automatizables | Media | `junior-developer`, `automation-specialist` |
| Validaciones críticas | Alta | `qa-engineer`, `security-expert` |
| Documentación técnica | Media | `tech-writer`, `developer` |
| Optimización de rendimiento | Alta | `performance-engineer`, `senior-developer` |

### Ejemplo de Asignación RACI

```yaml
fase_1_diseño:
  responsible: architect
  accountable: tech-lead
  consulted: [api-design, security-analysis]
  informed: [product-manager]
  
fase_2_implementacion:
  responsible: backend-developer
  accountable: architect
  consulted: [code-generation, testing, database-design]
  informed: [qa-engineer]
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
  consumers: [qa-automation, deployment-manager]
```

## Output y Artefactos

| Artefacto | Ubicación | Formato | Validador | Obligatorio |
|-----------|-----------|---------|-----------|-------------|
| Plan técnico | `.claude/plans/{timestamp}-{command-name}.md` | Markdown | `plan-validator` | Sí (planning) |
| Diagrama de arquitectura | `.claude/diagrams/{id}.mmd` | Mermaid | - | No |
| Checklist de validación | `.claude/checklists/{id}.json` | JSON | `schema-validator` | Sí |
| Reporte de análisis | `.claude/reports/{id}.md` | Markdown | - | Sí (executable) |
| Log de ejecución | `.claude/logs/{command-name}-{date}.log` | Plain text | - | Sí |

## Rollback y Cancelación

Si el command falla o el usuario cancela durante la ejecución:

### Procedimiento de Rollback

1. **Detener agentes en curso**: Enviar señal de cancelación a todos los agentes activos
2. **Eliminar artefactos parciales**:
    - Borrar planes incompletos en `.claude/plans/`
    - Limpiar archivos temporales en `.claude/temp/`
3. **Restaurar estado previo**: Si se modificó contexto compartido, revertir a snapshot anterior
4. **Registrar cancelación**:
   ```
   .claude/logs/cancelled-{timestamp}.log
   ```
5. **Notificar dependencias**: Informar a commands/MCPs que dependían de este output

### Estados Finales Posibles

- `completed`: Ejecución exitosa
- `failed`: Error irrecuperable
- `cancelled`: Cancelado por usuario
- `partial`: Completado parcialmente (solo para meta commands)

## Reglas Críticas

- **No modificación de código**: Bajo ningún concepto este command puede crear, modificar o eliminar archivos de código
- **Selección obligatoria de agentes**: El command debe elegir explícitamente los agentes y skills adecuados usando el framework RACI
- **Planificación obligatoria**: Si hay impacto en código, debe generarse un plan en `.claude/plans/`
- **Separación de responsabilidades**: Commands orquestan, agentes ejecutan
- **Ejecución restringida**: Solo `plan-manager` puede ejecutar planes
- **Análisis previo obligatorio**: Ninguna acción sin validación previa
- **Versionado semántico**: Cambios en el command requieren actualizar la versión
- **Idempotencia**: Ejecutar el command múltiples veces con los mismos parámetros debe producir el mismo resultado

---

## Acción del Usuario

{Prompt final claro y accionable.
Debe guiar al usuario sobre qué información proporcionar.

Ejemplo:
"Describe la feature que deseas implementar, incluyendo:
- Funcionalidad deseada
- Restricciones técnicas
- Criterios de aceptación
- Prioridad y timeline"}
```

---

## 4. Convenciones Obligatorias

### 4.1 Nombres de Commands

**Formato**: `{verbo}-{sustantivo}` (kebab-case)

✅ **Ejemplos válidos**:
- `analyze-feature`
- `plan-refactor`
- `audit-security`
- `generate-report`

❌ **Ejemplos inválidos**:
- `FeaturePlanner` (PascalCase)
- `analyze_code` (snake_case)
- `doEverything` (camelCase)
- `PLAN` (solo verbo)

### 4.2 Estructura de Directorios

```
.claude/
├── commands/
│   ├── planning/
│   │   ├── feature-planner.md
│   │   ├── refactor-analyzer.md
│   │   └── migration-coordinator.md
│   ├── executable/
│   │   ├── lint-runner.md
│   │   ├── security-audit.md
│   │   └── dependency-checker.md
│   └── meta/
│       └── plan-manager.md
├── plans/
│   ├── 20250120-143022-feature-planner.md
│   └── 20250120-150033-refactor-analyzer.md
├── logs/
│   ├── feature-planner-2025-01-20.log
│   └── cancelled-20250120-143555.log
├── reports/
│   └── security-audit-20250120.md
├── diagrams/
│   └── architecture-oauth2.mmd
├── checklists/
│   └── feature-validation.json
└── context/
└── shared-state.json
```

---

## 5. Anti-patrones Comunes

### ❌ God Command

**Problema**: Un command intenta hacer demasiado (análisis + planificación + ejecución + validación).

**Ejemplo**:
```markdown
name: do-everything
description: Analiza, planifica, ejecuta y valida cualquier tarea
```

**Solución**: Dividir en commands especializados:
- `analyze-requirements` (executable)
- `plan-implementation` (planning)
- `validate-output` (executable)

---

### ❌ Agent Micromanagement

**Problema**: El command especifica línea por línea qué debe hacer el agente, eliminando su autonomía.

**Ejemplo**:
```markdown
### Fase 1
- Crear variable `authToken` de tipo string
- Inicializarla en null
- Crear función `validateToken(token: string): boolean`
- Implementar lógica: if (token.length > 0) return true
```

**Solución**: Delegar la implementación, solo definir requisitos:
```markdown
### Fase 1: Gestión de Tokens (backend-developer)
- Implementar sistema de validación de tokens JWT
- Criterios: soporte RS256, expiración configurable, refresh tokens
- Skills: `auth-design`, `code-generation`
```

---

### ❌ Circular Dependencies

**Problema**: Command A invoca B, que invoca A, creando un loop infinito.

**Ejemplo**:
```
refactor-planner → code-analyzer → quality-checker → refactor-planner
```

**Solución**: Detectar ciclos en el análisis previo:
```json
{
  "validation_passed": false,
  "blocking_issues": [
    "Dependencia circular detectada: refactor-planner -> code-analyzer -> refactor-planner"
  ]
}
```

---

### ❌ Plan Without Context

**Problema**: Generar un plan sin solicitar suficiente información al usuario.

**Ejemplo**: Usuario dice "mejora la app" y el command genera un plan genérico de 50 pasos.

**Solución**: Implementar la sección "Contexto Requerido del Usuario" y validar que esté completa antes de proceder.

---

### ❌ Silent Failures

**Problema**: El command falla pero no registra logs ni notifica al usuario.

**Solución**: Todo fallo debe:
1. Escribir en `.claude/logs/`
2. Retornar JSON con error detallado
3. Ejecutar procedimiento de rollback
4. Notificar a commands dependientes

---

## 6. Criterios de Aceptación (Checklist de PR)

Un command está completo y listo para producción si cumple:

### Obligatorios (7/7)

- [ ] **Análisis inicial implementado** con validaciones JSON
- [ ] **Al menos 1 agente seleccionado explícitamente** con framework RACI
- [ ] **Flujo de trabajo** con fases numeradas y criterios de salida
- [ ] **Reglas críticas** documentadas y verificables
- [ ] **Ejemplo de uso** en la sección final con caso real
- [ ] **Documentación de rollback** con procedimiento paso a paso
- [ ] **Versionado semántico** en frontmatter

### Recomendados (5/5)

- [ ] Diagrama de flujo en Mermaid (`.claude/diagrams/`)
- [ ] Tests de validación automatizados
- [ ] Métricas de rendimiento esperadas
- [ ] Documentación de MCPs utilizados
- [ ] Ejemplos de output para cada tipo de resultado (éxito/fallo/cancelación)

**Calidad mínima para merge**: 7/7 obligatorios ✅  
**Calidad recomendada**: 12/12 (obligatorios + recomendados) ✅

---

## 7. Anexo A: Ejemplo Real Completo

### Command: `api-feature-planner`

```markdown
---
name: api-feature-planner
version: 1.0.0
author: platform-team
description: Analiza una nueva feature de API REST y genera un plan técnico de implementación
usage: "api-feature-planner [feature-description] [--priority=high]"
type: planning
writes_code: false
creates_plan: true
requires_approval: true
dependencies: [security-audit, api-design-validator]
---

# Comando: API Feature Planner

## Objetivo

Analizar una solicitud de nueva funcionalidad para una API REST y generar un plan técnico detallado que incluya:
- Diseño de endpoints
- Validaciones de seguridad
- Estrategia de testing
- Plan de migración si es necesario

**No ejecuta código**, solo coordina el análisis y genera la estrategia.

## Contexto Requerido del Usuario

- [ ] Descripción funcional de la feature (qué debe hacer)
- [ ] Endpoints involucrados (nuevos o modificados)
- [ ] Payload de ejemplo (request/response)
- [ ] Restricciones de seguridad (autenticación, autorización)
- [ ] SLA esperado (latencia, throughput)
- [ ] Versión de la API afectada (v1, v2, etc.)

## Análisis Inicial

### Validaciones Pre-ejecución

```json
{
  "validation_passed": true,
  "risks": [
    "Modifica esquema de base de datos",
    "Requiere nuevo servicio de autenticación"
  ],
  "required_approvals": ["tech-lead", "security-team"],
  "estimated_complexity": "high",
  "blocking_issues": []
}
```

## Selección de Agentes y Skills

### Fase 1: Diseño de Seguridad

```yaml
responsible: security-expert
accountable: architect
consulted: [security-analysis, threat-modeling, oauth-design]
informed: [compliance-team]
```

### Fase 2: Diseño de API

```yaml
responsible: api-architect
accountable: tech-lead
consulted: [api-design, openapi-generation, versioning-strategy]
informed: [frontend-team, mobile-team]
```

### Fase 3: Implementación

```yaml
responsible: backend-developer
accountable: senior-developer
consulted: [code-generation, testing, database-design]
informed: [qa-engineer, devops-team]
```

## Flujo de Trabajo Orquestado

### 1. Análisis de Seguridad (security-expert | Validado por architect)

**Objetivo**: Definir requisitos de autenticación y autorización

**Tareas**:
- Evaluar si la feature requiere OAuth2, API Keys o JWT
- Identificar datos sensibles en el payload
- Definir rate limiting necesario
- Documentar posibles vectores de ataque

**Asignación**:
- **Agente**: security-expert
- **Skills**: `security-analysis`, `threat-modeling`, `oauth-design`
- **MCPs**: `owasp-validator`
- **Validador**: architect

**Criterios de Salida**:
- [ ] Documento de análisis de amenazas generado
- [ ] Estrategia de autenticación definida
- [ ] Rate limits especificados

---

### 2. Diseño de Endpoints (api-architect | Validado por tech-lead)

**Objetivo**: Definir la estructura de los endpoints y contratos de datos

**Tareas**:
- Diseñar URIs según convenciones RESTful
- Definir esquemas JSON (request/response)
- Generar especificación OpenAPI 3.0
- Validar versionado de API

**Asignación**:
- **Agente**: api-architect
- **Skills**: `api-design`, `openapi-generation`, `versioning-strategy`
- **Dependencias**: Fase 1 completada
- **Validador**: tech-lead

**Criterios de Salida**:
- [ ] Especificación OpenAPI generada y validada
- [ ] Endpoints documentados con ejemplos
- [ ] Estrategia de versionado aprobada

---

### 3. Implementación (backend-developer | Validado por senior-developer)

**Objetivo**: Generar el código base de los endpoints

**Tareas**:
- Implementar controllers con validación de input
- Crear servicios de lógica de negocio
- Implementar capa de persistencia (si aplica)
- Escribir tests unitarios y de integración

**Asignación**:
- **Agente**: backend-developer
- **Skills**: `code-generation`, `testing`, `database-design`
- **Dependencias**: Fase 2 completada
- **Validador**: senior-developer

**Criterios de Salida**:
- [ ] Código implementado con coverage >80%
- [ ] Tests de integración pasando
- [ ] Documentación técnica actualizada

## Uso de otros Commands y MCPs

```yaml
commands_invocados:
  - name: security-audit
    trigger: post-fase-1
    output_required: security-report.json
    
  - name: api-design-validator
    trigger: post-fase-2
    output_required: openapi-validation.json

mcps_utilizados:
  - name: owasp-validator
    config: .claude/mcp-configs/owasp.json
    purpose: Validar contra top 10 de OWASP
    
  - name: database-schema-validator
    config: .claude/mcp-configs/db-validator.json
    purpose: Verificar migraciones compatibles
```

## Output y Artefactos

| Artefacto | Ubicación | Formato | Validador | Obligatorio |
|-----------|-----------|---------|-----------|-------------|
| Plan técnico | `.claude/plans/{timestamp}-api-feature.md` | Markdown | `plan-validator` | Sí |
| Especificación OpenAPI | `.claude/specs/api-v2-{feature}.yaml` | YAML | `openapi-validator` | Sí |
| Diagrama de secuencia | `.claude/diagrams/api-flow.mmd` | Mermaid | - | No |
| Checklist de seguridad | `.claude/checklists/security-{feature}.json` | JSON | `schema-validator` | Sí |

## Rollback y Cancelación

1. Eliminar plan parcial en `.claude/plans/`
2. Borrar especificaciones OpenAPI temporales
3. Notificar a `security-audit` y `api-design-validator`
4. Registrar en `.claude/logs/cancelled-{timestamp}.log`

## Reglas Críticas

- **No modificación de código**: Este command solo genera planes
- **Selección RACI obligatoria**: Cada fase debe tener responsible/accountable
- **Validación de seguridad**: Fase 1 es bloqueante
- **Aprobación requerida**: Plan debe ser aprobado por tech-lead antes de ejecución
- **Versionado de API**: Nunca modificar versiones existentes sin estrategia de deprecación

---

## Acción del Usuario

Describe la feature de API que deseas implementar, incluyendo:

1. **Funcionalidad**: ¿Qué debe hacer la API? (ej: "Autenticación OAuth2 para usuarios externos")
2. **Endpoints**: ¿Qué URIs necesitas? (ej: `POST /api/v2/auth/login`)
3. **Payload**: Proporciona ejemplos de request/response
4. **Seguridad**: ¿Qué nivel de protección necesita? (pública, autenticada, admin)
5. **SLA**: ¿Requisitos de rendimiento? (ej: "<200ms p95")
6. **Prioridad**: ¿Urgencia? (crítica, alta, media, baja)

**Ejemplo de solicitud válida**:
> "Necesito implementar autenticación OAuth2 para permitir que aplicaciones de terceros accedan a nuestra API. Endpoints: POST /api/v2/oauth/authorize, POST /api/v2/oauth/token. Debe soportar Authorization Code Grant. SLA: <500ms p99. Prioridad: alta."
```

---

## 8. Resumen de Reglas Inquebrantables

🔒 **Nunca tocar código** - Commands orquestan, no implementan  
🧠 **Analizar siempre antes de actuar** - Validación obligatoria pre-ejecución  
🧑‍💼 **Elegir explícitamente agentes y skills** - Framework RACI en todas las fases  
🗂 **Planes solo en `.claude/plans/`** - Ubicación estandarizada  
🧭 **Commands orquestan, agentes ejecutan** - Separación de responsabilidades  
🗝 **`plan-manager` único ejecutor** - Centralización de ejecución de planes  
📊 **Versionado semántico** - Cambios rastreables  
🔄 **Idempotencia garantizada** - Mismos inputs = mismo output  
🚨 **Rollback documentado** - Procedimiento de cancelación obligatorio  
✅ **7/7 criterios de calidad** - No merge sin completar checklist

---

## 9. Recursos Adicionales

### Plantillas Disponibles

- `.claude/templates/planning-command.md` - Plantilla para planning commands
- `.claude/templates/executable-command.md` - Plantilla para executable commands
- `.claude/templates/meta-command.md` - Plantilla para meta commands

### Validadores

- `plan-validator` - Valida estructura de planes generados
- `schema-validator` - Valida JSON contra esquemas definidos
- `agent-selector-validator` - Verifica asignaciones RACI

### Comandos de Utilidad

```bash
# Validar un command antes de commit
claude validate-command ./commands/planning/my-new-command.md

# Generar command desde template
claude generate-command --type planning --name feature-planner

# Verificar dependencias circulares
claude check-circular-deps
```

---

**Versión de la guía**: 2.0.0  
**Última actualización**: 2025-01-20
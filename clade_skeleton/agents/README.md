# Guía Técnica para la Creación de Agentes en Claude Code

## 1. Principio Fundamental: El Agente No Sabe, El Agente Razona

Un agente es una **entidad deliberadamente no especializada** cuyo propósito no es poseer conocimiento técnico concreto, sino **aplicar razonamiento estructurado** sobre un contexto dado.

### El agente define:
- ✅ Cómo analiza problemas
- ✅ Cómo toma decisiones
- ✅ Cómo evalúa resultados
- ✅ Cómo itera hasta cumplir un objetivo

### El agente NO define:
- ❌ Lenguajes de programación
- ❌ Frameworks específicos
- ❌ Herramientas concretas
- ❌ Entornos de ejecución

**Todo el conocimiento técnico son dependencias externas inyectadas en tiempo de ejecución.**

---

## 2. Arquitectura de Tres Capas

### 2.1 Agente (Razonamiento)
- Rol abstracto (ej: "Senior Backend Engineer")
- Principios generales (SOLID, diseño orientado a mantenibilidad)
- Capacidad de planificación y verificación
- Orquestación del bucle de ejecución

👉 **El agente es agnóstico al dominio técnico**

### 2.2 Skills (Conocimiento)
- Conocimiento técnico declarativo
- Lenguajes, frameworks, estilos arquitectónicos
- Convenciones de equipo y proyecto
- **Se inyectan en tiempo de invocación**

👉 **El agente no es experto en nada por defecto**

### 2.3 Tools (Capacidad de Acción)
- Habilidades operativas
- Ejecutar comandos, leer/escribir archivos
- Consultar APIs, acceder a sistemas externos
- **"Darle un ordenador al agente"**

👉 **Sin tools, el agente solo razona. Con tools, el agente actúa.**

---

## 3. El Bucle Operativo (Agent Loop)

Todo agente opera en un ciclo explícito y controlado:

```
┌─────────────────────────────────────────┐
│  1. RECOPILAR CONTEXTO                  │
│  - Leer archivos                        │
│  - Consultar estado del sistema         │
│  - Revisar resultados previos           │
│  👉 El contexto no se asume, se obtiene │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  2. TOMAR ACCIÓN                        │
│  - Ejecutar comandos                    │
│  - Modificar archivos                   │
│  - Aplicar skills inyectadas            │
│  👉 Acciones solo vía tools explícitas  │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  3. VERIFICAR RESULTADOS                │
│  - Validar salidas                      │
│  - Comprobar errores                    │
│  - Comparar con objetivo esperado       │
│  👉 No confiar, verificar               │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  4. DECIDIR                             │
│  ¿Objetivo cumplido?                    │
│   SI → Finalizar                        │
│   NO → Ajustar plan y volver a (1)     │
│  👉 Iteración hasta éxito o límite      │
└─────────────────────────────────────────┘
```

---

## 4. Template Oficial de Agentes

```markdown
---
name: {agent-name}
version: 1.0.0
author: {team/person}
description: {Rol abstracto del agente. Enfocado en razonamiento, NO en tecnologías.}
model: {model-id}
color: {color-hex}
type: {reasoning | validation | orchestration}
autonomy_level: {low | medium | high}
requires_human_approval: {true | false}
max_iterations: {number}
---

# Agente: {Agent Title}

## 1. Perfil de Razonamiento

### 1.1 Identidad Abstracta
- **Rol**: {Ej: Arquitecto de Sistemas / Especialista en Seguridad / QA Engineer}
- **Mentalidad**: {Ej: Defensiva (seguridad) / Pragmática (entrega) / Optimizadora (rendimiento)}
- **Alcance de Responsabilidad**: {Ej: Backend APIs / Frontend Components / Infrastructure}

### 1.2 Principios de Diseño
Estos principios guían cada decisión del agente:

- {Principio 1}: {Ej: SOLID - Código debe ser extensible sin modificación}
- {Principio 2}: {Ej: KISS - Preferir soluciones simples sobre complejas}
- {Principio 3}: {Ej: Security by Design - Validar inputs, sanitizar outputs}
- {Principio 4}: {Ej: Fail Fast - Detectar errores lo antes posible}

### 1.3 Objetivo Final
{Descripción clara del resultado esperado cuando el agente completa su tarea.}

**Ejemplo**: Garantizar que todo código entregado:
- Cumple con los tests definidos
- Sigue las convenciones del proyecto
- Está documentado adecuadamente
- No introduce regresiones

---

## 2. Bucle Operativo (Agent Loop)

Este agente opera bajo un ciclo estrictamente controlado. **Cada iteración debe ser verificable y auditable.**

### 2.1 Fase: RECOPILAR CONTEXTO

**Regla de Oro**: No asumir estados previos. Todo debe ser verificado empíricamente.

**Acciones permitidas**:
- Leer archivos del proyecto (configs, código existente)
- Consultar logs de ejecuciones previas
- Inspeccionar estado del sistema (git status, procesos, etc.)
- Revisar outputs de tools previas

**Output esperado**:
```json
{
  "context_gathered": true,
  "files_read": ["src/main.ts", "package.json"],
  "system_state": {
    "git_branch": "feature/oauth",
    "uncommitted_changes": false
  },
  "previous_errors": []
}
```

---

### 2.2 Fase: PLANIFICACIÓN Y ACCIÓN

**Regla de Oro**: Aplicar skills inyectadas + ejecutar vía tools explícitas.

**Proceso de decisión**:
1. Identificar qué skills son relevantes para la tarea actual
2. Formular un plan de acción basado en el conocimiento inyectado
3. Seleccionar las tools necesarias
4. Ejecutar acciones paso a paso
5. Registrar cada acción en logs

**Ejemplo de razonamiento**:
```
Tarea: Implementar endpoint POST /api/users

Skills disponibles: [TypeScriptSkill, ExpressSkill, CleanArchitectureSkill]
Tools disponibles: [FileSystem, Terminal, TestRunner]

Plan:
1. [TypeScript + Clean Architecture] Crear UserController en src/controllers/
2. [Express] Registrar ruta en src/routes/users.routes.ts
3. [FileSystem] Escribir código en archivos
4. [Terminal] Ejecutar npm run build para verificar compilación
5. [TestRunner] Ejecutar tests de integración
```

**Output esperado**:
```json
{
  "plan_executed": true,
  "actions_taken": [
    {
      "tool": "FileSystem",
      "action": "write",
      "file": "src/controllers/UserController.ts",
      "success": true
    },
    {
      "tool": "Terminal",
      "command": "npm run build",
      "exit_code": 0
    }
  ]
}
```

---

### 2.3 Fase: VERIFICACIÓN

**Regla de Oro**: No confiar en que algo funcionó. Verificarlo explícitamente.

**Checklist de verificación** (adaptar según tipo de tarea):
- [ ] ¿El código compila sin errores?
- [ ] ¿Los tests pasan?
- [ ] ¿Se siguen las convenciones del proyecto?
- [ ] ¿No hay warnings críticos?
- [ ] ¿El resultado coincide con el objetivo esperado?

**Métodos de verificación**:
```yaml
compilacion:
  tool: Terminal
  command: "npm run build"
  success_criteria: "exit_code == 0"

tests:
  tool: TestRunner
  command: "npm test -- --coverage"
  success_criteria: "all_passed && coverage > 80%"

linting:
  tool: Terminal
  command: "npm run lint"
  success_criteria: "exit_code == 0"
```

**Output esperado**:
```json
{
  "verification_passed": true,
  "checks_performed": [
    {"name": "compilation", "passed": true},
    {"name": "tests", "passed": true, "coverage": 85},
    {"name": "linting", "passed": true}
  ],
  "issues_found": []
}
```

---

### 2.4 Fase: ITERACIÓN

**Regla de Oro**: Ajustar el plan basándose en resultados empíricos.

**Criterios de decisión**:
```
SI (verificación exitosa) Y (objetivo cumplido):
    → FINALIZAR con éxito

SI (verificación exitosa) Y (objetivo parcialmente cumplido):
    → CONTINUAR con siguiente sub-tarea

SI (verificación fallida) Y (iteraciones < max_iterations):
    → ANALIZAR error
    → AJUSTAR plan
    → VOLVER a fase de acción

SI (iteraciones >= max_iterations):
    → ESCALAR a humano
    → REPORTAR estado y errores
```

**Output de iteración**:
```json
{
  "iteration": 3,
  "status": "retrying",
  "reason": "Tests fallaron - error en validación de email",
  "adjustment": "Agregar regex de validación en UserValidator",
  "next_action": "modificar src/validators/UserValidator.ts"
}
```

---

## 3. Capacidades Inyectadas (Runtime Configuration)

**IMPORTANTE**: Este agente **no posee conocimiento técnico intrínseco**. Su efectividad depende de los recursos proporcionados en la invocación.

### 3.1 Skills (Conocimiento Declarativo)

Las skills se inyectan como contexto estructurado:

```typescript
interface Skill {
  name: string;
  version: string;
  description: string;
  conventions: string[];
  best_practices: string[];
  anti_patterns: string[];
  examples: CodeExample[];
}
```

**Ejemplo de inyección**:
```json
{
  "skills": [
    {
      "name": "TypeScriptSkill",
      "version": "5.0",
      "conventions": [
        "Usar tipos explícitos, evitar any",
        "Interfaces para contratos públicos",
        "Types para uniones y utilidades"
      ],
      "best_practices": [
        "Preferir unknown sobre any para inputs no validados",
        "Usar strict mode en tsconfig.json"
      ],
      "anti_patterns": [
        "Usar ! (non-null assertion) sin justificación",
        "Type casting con 'as' sin validación previa"
      ]
    },
    {
      "name": "CleanArchitectureSkill",
      "version": "1.0",
      "conventions": [
        "Estructura: controllers -> services -> repositories",
        "Dependency injection mediante interfaces",
        "Separar lógica de negocio de infraestructura"
      ]
    }
  ]
}
```

**Aplicación en el agente**:
El agente consulta las skills antes de cada decisión técnica y las aplica como restricciones.

---

### 3.2 Tools (Capacidad de Acción)

Las tools otorgan al agente "acceso al ordenador":

```typescript
interface Tool {
  name: string;
  capabilities: string[];
  permissions: Permission[];
  rate_limits?: RateLimit;
}
```

**Ejemplo de configuración**:
```yaml
tools:
  - name: FileSystem
    capabilities:
      - read_file
      - write_file
      - list_directory
      - create_directory
    permissions:
      allowed_paths: ["src/", "tests/", "docs/"]
      forbidden_paths: [".env", "node_modules/", ".git/"]
      max_file_size: 1MB
    
  - name: Terminal
    capabilities:
      - execute_command
      - read_stdout
      - read_stderr
    permissions:
      allowed_commands: ["npm", "git", "tsc", "jest"]
      forbidden_commands: ["rm -rf", "sudo", ":(){:|:&};:"]
      timeout: 30s
    
  - name: TestRunner
    capabilities:
      - run_unit_tests
      - run_integration_tests
      - generate_coverage
    permissions:
      test_frameworks: ["jest", "vitest"]
      
  - name: APIClient
    capabilities:
      - http_get
      - http_post
    permissions:
      allowed_domains: ["api.internal.company.com"]
      require_auth: true
```

**Restricciones críticas**:
- Agente solo puede usar tools explícitamente inyectadas
- Toda acción debe pasar por una tool (no alucinaciones)
- Permisos de tools son inmutables durante ejecución

---

## 4. Estrategia de Toma de Decisiones

Define el **modelo mental** que el agente debe seguir al enfrentarse a decisiones.

### 4.1 Análisis de Impacto

Antes de modificar código, el agente debe evaluar:

**Framework de evaluación**:
```
Cambio Propuesto: {descripción}

Impacto en:
├── Arquitectura: {bajo | medio | alto}
├── Seguridad: {bajo | medio | alto}
├── Rendimiento: {bajo | medio | alto}
├── Mantenibilidad: {mejor | neutral | peor}
└── Breaking Changes: {sí | no}

Decisión:
SI (algún impacto == alto) O (breaking_changes == sí):
    → Generar plan y solicitar aprobación humana
SINO:
    → Proceder con la implementación
```

---

### 4.2 Priorización de Tareas

Cuando hay múltiples sub-tareas, el agente debe seguir este orden:

1. **Crítico (bloqueantes)**: Errores de compilación, tests rotos
2. **Alto (seguridad)**: Validaciones, sanitización, autenticación
3. **Medio (funcionalidad)**: Implementación de features
4. **Bajo (mejoras)**: Refactoring, optimizaciones

**Ejemplo**:
```
Tareas pendientes:
- [CRÍTICO] Fix: Endpoint /api/users retorna 500
- [ALTO] Agregar validación de JWT en middleware
- [MEDIO] Implementar paginación en /api/posts
- [BAJO] Refactor: Extraer lógica duplicada en utils

Orden de ejecución: CRÍTICO → ALTO → MEDIO → BAJO
```

---

### 4.3 Gestión de Errores

Define **estrategias específicas** para errores comunes:

```yaml
error_strategies:
  - error_type: "TypeScript compilation error"
    strategy: |
      1. Leer mensaje de error completo
      2. Localizar archivo y línea afectada
      3. Consultar TypeScriptSkill para convenciones
      4. Aplicar fix siguiendo convenciones
      5. Re-compilar y verificar
      6. Si persiste después de 3 intentos → Escalar
      
  - error_type: "Test failure"
    strategy: |
      1. Identificar test fallido y assertion
      2. Ejecutar solo ese test con --verbose
      3. Revisar código bajo test
      4. Aplicar fix según lógica del test
      5. Re-ejecutar suite completa
      6. Si coverage baja → Agregar tests faltantes
      
  - error_type: "Linting error"
    strategy: |
      1. Ejecutar linter con --fix si disponible
      2. Si no se auto-corrige, leer regla violada
      3. Aplicar corrección manual
      4. Re-ejecutar linter
      5. Si regla es cuestionable → Documentar y notificar
```

---

### 4.4 Escalación a Humanos

El agente debe **reconocer sus límites** y escalar cuando:

- ❌ Después de `max_iterations` sin éxito
- ❌ Cambio requiere decisión arquitectónica mayor
- ❌ Herramienta necesaria no está disponible
- ❌ Contexto insuficiente para continuar
- ❌ Conflicto entre skills (convenciones contradictorias)

**Formato de escalación**:
```json
{
  "escalation_reason": "unable_to_resolve_after_max_iterations",
  "iterations_completed": 5,
  "last_error": "Test 'UserController.createUser' fails with 'Email validation error'",
  "attempted_solutions": [
    "Added regex validation in UserValidator",
    "Updated email schema in Joi",
    "Fixed typo in validation logic"
  ],
  "context_provided": {
    "files_modified": ["src/validators/UserValidator.ts"],
    "logs": ".claude/logs/backend-engineer-2025-01-20.log"
  },
  "recommended_next_steps": "Review email validation requirements with Product team"
}
```

---

## 5. Reglas de Oro (Invariantes del Agente)

Estas reglas **nunca** deben violarse:

### 5.1 No Alucinar
- ❌ **NUNCA** asumir que un comando funcionó sin verificarlo
- ❌ **NUNCA** inventar paths de archivos que no existen
- ❌ **NUNCA** afirmar conocimiento técnico que no está en las skills inyectadas

✅ **SIEMPRE** verificar con tools antes de afirmar

---

### 5.2 Verificación Empírica
- ❌ Confiar en que `npm run build` funcionó por "lógica"
- ✅ Ejecutar `npm run build` y verificar `exit_code === 0`

---

### 5.3 Trazabilidad
Todo cambio significativo debe:
1. Registrarse en `.claude/logs/{agent-name}-{date}.log`
2. Incluir razonamiento: "¿Por qué este cambio?"
3. Referenciar skill aplicada: "Según CleanArchitectureSkill..."

**Ejemplo de log**:
```
[2025-01-20 14:30:22] backend-engineer
ACCIÓN: Crear archivo src/controllers/UserController.ts
RAZÓN: Implementar endpoint POST /api/users según plan
SKILL APLICADA: CleanArchitectureSkill - separación de concerns
VERIFICACIÓN: Compilación exitosa, 0 errores
```

---

### 5.4 Idempotencia
Ejecutar el agente múltiples veces con el mismo input debe:
- Producir el mismo resultado
- No causar efectos secundarios no deseados

---

### 5.5 Fail-Safe Defaults
Ante ambigüedad, el agente debe:
- ❌ **NO** elegir la opción "más avanzada"
- ✅ **SÍ** elegir la opción **más simple y segura**

**Ejemplo**: Si no está claro si usar `any` o `unknown`:
```typescript
// ❌ NO hacer por defecto
function process(data: any) { ... }

// ✅ SÍ hacer por defecto (más seguro)
function process(data: unknown) { ... }
```

---

## 6. Restricciones y Políticas

### 6.1 Seguridad

```yaml
security_policies:
  - rule: "No leer archivos fuera de allowed_paths"
    enforcement: "FileSystem tool rechaza acceso"
    
  - rule: "No ejecutar comandos no whitelisteados"
    enforcement: "Terminal tool bloquea ejecución"
    
  - rule: "No exponer secrets en logs"
    enforcement: "Logger sanitiza valores sensibles automáticamente"
    
  - rule: "Validar inputs antes de uso"
    enforcement: "Skill de seguridad requiere validación explícita"
```

---

### 6.2 Entorno

```yaml
environment_rules:
  - rule: "Ejecutar tests antes de marcar tarea como completa"
    verification: "TestRunner tool debe retornar all_passed: true"
    
  - rule: "No hacer commit sin linter pasando"
    verification: "npm run lint debe retornar exit_code 0"
    
  - rule: "Documentar funciones públicas"
    verification: "Verificar JSDoc en exports de módulos"
```

---

### 6.3 Límites Operacionales

```yaml
operational_limits:
  max_iterations: 10
  max_file_size: 1MB
  max_execution_time: 5m
  max_parallel_tools: 3
  
  on_limit_exceeded:
    action: "escalate_to_human"
    include: ["logs", "context", "attempted_solutions"]
```

---

## 7. Tipos de Agentes

### 7.1 Reasoning Agents (Razonamiento)
**Propósito**: Análisis, diseño, planificación

**Características**:
- Alto uso de skills
- Bajo uso de tools (solo lectura)
- Output: Planes, diagramas, documentación

**Ejemplo**: `architect-agent`, `design-reviewer`

---

### 7.2 Validation Agents (Validación)
**Propósito**: QA, testing, seguridad, auditoría

**Características**:
- Medio uso de skills (conocimiento de buenas prácticas)
- Medio uso de tools (ejecutar tests, linters)
- Output: Reportes de validación, checklists

**Ejemplo**: `qa-engineer`, `security-auditor`

---

### 7.3 Orchestration Agents (Orquestación)
**Propósito**: Coordinación de múltiples agentes/tasks

**Características**:
- Bajo uso de skills (generalistas)
- Alto uso de tools (ejecutar, monitorear)
- Output: Flujos de trabajo, estados de ejecución

**Ejemplo**: `workflow-coordinator`, `release-manager`

---

## 8. Ejemplo Real Completo

```markdown
---
name: backend-engineer
version: 1.0.0
author: platform-team
description: Senior Backend Engineer especializado en razonamiento sobre APIs REST y servicios backend
model: claude-sonnet-4-20250514
color: "#3B82F6"
type: reasoning
autonomy_level: medium
requires_human_approval: false
max_iterations: 10
---

# Agente: Backend Engineer

## 1. Perfil de Razonamiento

### 1.1 Identidad Abstracta
- **Rol**: Senior Backend Engineer
- **Mentalidad**: Pragmática - equilibrio entre calidad y entrega
- **Alcance de Responsabilidad**: APIs REST, servicios backend, integraciones

### 1.2 Principios de Diseño
- **SOLID**: Código debe ser extensible sin modificación (Open/Closed)
- **KISS**: Preferir soluciones simples sobre ingeniería excesiva
- **Fail Fast**: Validar inputs en el borde del sistema, fallar temprano
- **Separation of Concerns**: Controllers, Services, Repositories claramente separados

### 1.3 Objetivo Final
Entregar código backend que:
- Pasa todos los tests (unit + integration)
- Sigue las convenciones del proyecto
- Tiene cobertura > 80%
- No introduce vulnerabilidades de seguridad
- Está documentado con JSDoc/comentarios donde es complejo

---

## 2. Bucle Operativo

### 2.1 RECOPILAR CONTEXTO

Acciones:
1. Leer package.json para entender stack (Express? Fastify? NestJS?)
2. Consultar tsconfig.json para configuración TypeScript
3. Revisar estructura de directorios (src/controllers, src/services, etc.)
4. Leer .eslintrc y .prettierrc para convenciones de código
5. Consultar tests existentes para entender patrones de testing

Output esperado:
```json
{
  "context": {
    "framework": "Express 4.18",
    "typescript_version": "5.0",
    "test_framework": "Jest",
    "architecture": "Clean Architecture (3 layers)"
  }
}
```

### 2.2 PLANIFICACIÓN Y ACCIÓN

Para tarea: "Implementar POST /api/users"

Plan:
1. **[CleanArchitectureSkill]** Identificar capa: Controller
2. **[TypeScriptSkill]** Crear interfaz `CreateUserDTO`
3. **[ExpressSkill]** Implementar route handler
4. **[FileSystem]** Escribir src/controllers/UserController.ts
5. **[Terminal]** Ejecutar `npm run build`
6. **[TestRunner]** Ejecutar `npm test`

### 2.3 VERIFICACIÓN

Checklist:
- [ ] Compilación: `tsc --noEmit` retorna exit 0
- [ ] Tests: `npm test` retorna all passed
- [ ] Linting: `npm run lint` retorna exit 0
- [ ] Coverage: > 80% en nuevo código

### 2.4 ITERACIÓN

```
SI (todos los checks pasan):
    → FINALIZAR
SI (algún check falla) Y (iteration < 10):
    → Analizar error específico
    → Aplicar fix según error_strategies
    → REPETIR desde 2.2
SI (iteration >= 10):
    → ESCALAR con contexto completo
```

---

## 3. Capacidades Inyectadas

### 3.1 Skills Esperadas
```json
{
  "required": ["TypeScriptSkill", "NodeSkill"],
  "optional": ["ExpressSkill", "NestJSSkill", "FastifySkill"],
  "architecture": ["CleanArchitectureSkill", "HexagonalArchitectureSkill"]
}
```

### 3.2 Tools Necesarias
```yaml
- FileSystem:
    permissions:
      read: ["src/", "tests/", "package.json", "tsconfig.json"]
      write: ["src/controllers/", "src/services/", "tests/"]
      
- Terminal:
    allowed_commands: ["npm", "tsc", "jest", "git"]
    timeout: 60s
    
- TestRunner:
    frameworks: ["jest", "vitest"]
    
- APIClient:
    allowed_domains: ["localhost:3000"]
```

---

## 4. Estrategia de Toma de Decisiones

### 4.1 Análisis de Impacto
```
Cambio: Agregar nuevo endpoint

Evaluación:
- Arquitectura: BAJO (sigue patrón existente)
- Seguridad: MEDIO (requiere validación de input)
- Rendimiento: BAJO (CRUD simple)
- Breaking Changes: NO

Decisión: PROCEDER sin aprobación
```

### 4.2 Priorización
1. CRÍTICO: Tests rotos, errores de compilación
2. ALTO: Validaciones de seguridad, autenticación
3. MEDIO: Features nuevas
4. BAJO: Refactoring, optimizaciones

### 4.3 Gestión de Errores
```yaml
- error: "TS2345: Argument of type 'string' is not assignable to parameter of type 'number'"
  strategy: |
    1. Leer TypeScriptSkill para convenciones de tipos
    2. Verificar si debe ser number o ajustar función
    3. Aplicar fix
    4. Re-compilar
    
- error: "Jest: Expected 201, received 500"
  strategy: |
    1. Ejecutar test con --verbose
    2. Revisar logs del servidor
    3. Identificar origen del 500 (validation? database?)
    4. Aplicar fix
    5. Re-ejecutar test
```

---

## 5. Reglas de Oro

- **No Alucinar**: Si no sé qué framework se usa, leer package.json antes de asumir
- **Verificación Empírica**: Ejecutar `npm run build` y verificar exit code, no confiar
- **Trazabilidad**: Registrar en logs por qué elegí Express middleware vs NestJS interceptor

---

## 6. Restricciones y Políticas

### Seguridad
- Validar todo input con Joi/Zod antes de procesarlo
- No exponer stack traces en producción
- Sanitizar outputs antes de enviar respuestas

### Entorno
- Tests obligatorios antes de marcar tarea completa
- Linter debe pasar antes de commit
- Documentar endpoints con JSDoc + OpenAPI

---

## 7. Invocación de Ejemplo

```typescript
await invokeAgent({
  agent: "backend-engineer",
  task: "Implementar POST /api/users con validación de email",
  skills: [
    TypeScriptSkill,
    ExpressSkill,
    CleanArchitectureSkill,
    JoiValidationSkill
  ],
  tools: [
    FileSystemTool,
    TerminalTool,
    TestRunnerTool
  ],
  constraints: {
    max_iterations: 10,
    required_coverage: 80,
    must_pass_linter: true
  }
});
```

**Output esperado**:
```json
{
  "status": "success",
  "iterations": 3,
  "files_modified": [
    "src/controllers/UserController.ts",
    "src/validators/UserValidator.ts",
    "tests/integration/users.test.ts"
  ],
  "verification": {
    "compilation": "passed",
    "tests": "passed (12/12)",
    "coverage": 87,
    "linting": "passed"
  }
}
```
```

---

## 9. Anti-patrones en Agentes

### ❌ Agente Omnisciente
**Problema**: Agente que "ya sabe" todo sin consultar skills

**Ejemplo**:
```markdown
## Conocimiento Intrínseco
- Experto en TypeScript, React, Node.js, PostgreSQL
- Conoce todas las mejores prácticas de seguridad
```

**Por qué es malo**: Viola el principio de inyección de dependencias. El agente no debe tener conocimiento hardcodeado.

**Solución**:
```markdown
## Capacidades Inyectadas
El agente aplicará las convenciones y frameworks definidos en las skills cargadas dinámicamente.
```

---

### ❌ Agente Sin Verificación
**Problema**: Confía en que las acciones funcionaron sin verificar

**Ejemplo**:
```markdown
1. Crear archivo controller.ts
2. Continuar con siguiente paso
```

**Por qué es malo**: Puede continuar con archivos no creados, generando cascada de errores.

**Solución**:
```markdown
1. Crear archivo controller.ts
2. [FileSystem] Verificar que archivo existe
3. [Terminal] Compilar y verificar exit_code
4. SI verificación exitosa → Continuar
```

---

### ❌ Agente Hardcodeado
**Problema**: Asume estructura de proyecto específica

**Ejemplo**:
```markdown
### Estructura esperada
- src/controllers/
- src/services/
- src/repositories/
```

**Por qué es malo**: Solo funciona con un tipo de proyecto.

**Solución**:
```markdown
### Fase: Descubrimiento
1. Leer estructura de directorios con FileSystem
2. Adaptar estrategia según arquitectura encontrada
3. Si no hay estructura clara → Sugerir organización
```

---

### ❌ Agente que Alucina Tools
**Problema**: Inventa comandos o tools que no existen

**Ejemplo**:
```markdown
1. Ejecutar `magic-deploy --auto`
2. Verificar con `check-deployment-status`
```

**Por qué es malo**: Estas tools no existen, el agente fallará silenciosamente.

**Solución**:
```markdown
### Tools Requeridas
- Terminal (con allowed_commands: ["npm", "git"])
- Si necesita desplegar → Requiere DeploymentTool explícita
- Si tool no disponible → Escalar a humano
```

---

### ❌ Loop Infinito Sin Max Iterations
**Problema**: Agente que puede iterar indefinidamente

**Ejemplo**:
```markdown
REPETIR hasta que tests pasen:
    - Modificar código
    - Ejecutar tests
```

**Por qué es malo**: Si hay un bug imposible de resolver automáticamente, el agente nunca termina.

**Solución**:
```markdown
---
max_iterations: 10
---

SI iteration >= max_iterations:
    → Escalar a humano con contexto completo
```

---

## 10. Convenciones de Nomenclatura

### Nombres de Agentes
**Formato**: `{rol}-{especialización}` (kebab-case)

✅ **Válidos**:
- `backend-engineer`
- `frontend-specialist`
- `qa-automation`
- `security-auditor`

❌ **Inválidos**:
- `BackendEngineer` (PascalCase)
- `backend_engineer` (snake_case)
- `engineer` (demasiado genérico)
- `do-everything-agent` (viola SRP)

---

## 11. Estructura de Directorios

```
.claude/
├── agents/
│   ├── reasoning/
│   │   ├── backend-engineer.md
│   │   ├── architect.md
│   │   └── design-reviewer.md
│   ├── validation/
│   │   ├── qa-engineer.md
│   │   ├── security-auditor.md
│   │   └── performance-analyzer.md
│   └── orchestration/
│       ├── workflow-coordinator.md
│       └── release-manager.md
├── skills/
│   ├── languages/
│   │   ├── typescript-skill.json
│   │   └── python-skill.json
│   └── frameworks/
│       ├── express-skill.json
│       └── react-skill.json
└── tools/
    ├── filesystem-tool.json
    ├── terminal-tool.json
    └── test-runner-tool.json
```

---

## 12. Criterios de Aceptación

Un agente está completo si cumple:

### Obligatorios (8/8)
- [ ] Perfil de razonamiento definido (rol + principios + objetivo)
- [ ] Bucle operativo completo (4 fases documentadas)
- [ ] Capacidades inyectadas especificadas (skills + tools)
- [ ] Estrategia de toma de decisiones con ejemplos
- [ ] Reglas de oro documentadas
- [ ] Restricciones y políticas explícitas
- [ ] Configuración de max_iterations y escalación
- [ ] Ejemplo de invocación con output esperado

### Recomendados (4/4)
- [ ] Diagrama de flujo del Agent Loop
- [ ] Anti-patrones específicos del dominio
- [ ] Métricas de éxito/fallo
- [ ] Tests de validación del agente

**Calidad mínima**: 8/8 obligatorios ✅  
**Calidad recomendada**: 12/12 ✅

---

## 13. Template Vacío Listo para Usar

```markdown
---
name: {agent-name}
version: 1.0.0
author: {team/person}
description: {Rol abstracto enfocado en razonamiento}
model: {model-id}
color: "{hex-color}"
type: {reasoning | validation | orchestration}
autonomy_level: {low | medium | high}
requires_human_approval: {true | false}
max_iterations: {number}
---

# Agente: {Agent Title}

## 1. Perfil de Razonamiento

### 1.1 Identidad Abstracta
- **Rol**: 
- **Mentalidad**: 
- **Alcance de Responsabilidad**: 

### 1.2 Principios de Diseño
- {Principio 1}: 
- {Principio 2}: 
- {Principio 3}: 

### 1.3 Objetivo Final


---

## 2. Bucle Operativo

### 2.1 RECOPILAR CONTEXTO

**Acciones**:
1. 
2. 
3. 

**Output esperado**:
```json
{
  "context_gathered": true
}
```

### 2.2 PLANIFICACIÓN Y ACCIÓN

**Proceso**:
1.
2.

**Output esperado**:
```json
{
  "plan_executed": true
}
```

### 2.3 VERIFICACIÓN

**Checklist**:
- [ ] 
- [ ] 

**Output esperado**:
```json
{
  "verification_passed": true
}
```

### 2.4 ITERACIÓN

```
SI (verificación exitosa):
    → FINALIZAR
SI (verificación fallida) Y (iteration < max):
    → Ajustar y reintentar
SI (iteration >= max):
    → Escalar
```

---

## 3. Capacidades Inyectadas

### 3.1 Skills Esperadas
```json
{
  "required": [],
  "optional": []
}
```

### 3.2 Tools Necesarias
```yaml
- ToolName:
    permissions: {}
```

---

## 4. Estrategia de Toma de Decisiones

### 4.1 Análisis de Impacto


### 4.2 Priorización


### 4.3 Gestión de Errores
```yaml
- error: ""
  strategy: |
```

---

## 5. Reglas de Oro

- **No Alucinar**:
- **Verificación Empírica**:
- **Trazabilidad**:

---

## 6. Restricciones y Políticas

### Seguridad


### Entorno


---

## 7. Invocación de Ejemplo

```typescript
await invokeAgent({
  agent: "{agent-name}",
  task: "",
  skills: [],
  tools: [],
  constraints: {}
});
```

**Output esperado**:
```json
{
  "status": "success"
}
```
```

---

## 14. Resumen Ejecutivo

### Principios Fundamentales

🧠 **Agente = Razonamiento puro** - Sin conocimiento técnico hardcodeado  
📚 **Skills = Conocimiento inyectado** - Convenciones, frameworks, lenguajes  
🛠️ **Tools = Capacidad de acción** - "Darle un ordenador al agente"  
🔄 **Loop = Autonomía controlada** - Observar → Actuar → Verificar → Repetir  
🚨 **Escalación = Reconocer límites** - Cuando max_iterations se agota  
✅ **Verificación = No confiar** - Todo debe comprobarse empíricamente  
📝 **Trazabilidad = Auditoría** - Logs de razonamiento y acciones  

---

## 15. Recursos Adicionales

### Plantillas
- `.claude/templates/reasoning-agent.md`
- `.claude/templates/validation-agent.md`
- `.claude/templates/orchestration-agent.md`

### Validadores
- `agent-validator` - Valida estructura del agente
- `loop-validator` - Verifica que el bucle está completo
- `skill-compatibility-checker` - Valida skills requeridas vs disponibles

### Comandos de Utilidad
```bash
# Validar agente
claude validate-agent ./agents/reasoning/backend-engineer.md

# Generar agente desde template
claude generate-agent --type reasoning --name backend-engineer

# Verificar compatibilidad skills
claude check-skills --agent backend-engineer --skills typescript,express
```

---

**Versión de la guía**: 2.0.0  
**Última actualización**: 2025-01-20  

# AGENTS.md

Este archivo proporciona la guía maestra para los agentes de IA (Claude Code, Cursor, Copilot, etc.) sobre cómo interactuar, crear y mantener las entidades del ecosistema Claude en este repositorio.

## 📁 Estructura del Ecosistema (.claude/)

La inteligencia y el flujo de trabajo del repositorio se centralizan en el directorio `.claude/`, utilizando las plantillas definidas en `clade_skeleton/` como base para la creación de nuevas entidades.

```
.claude/
  agents/               # Definiciones de agentes especializados
  commands/             # Flujos de trabajo orquestados (comandos)
  plans/                # Hojas de ruta de tareas
    active/             # Planes en ejecución
    completed/          # Histórico de planes finalizados
  sessions/             # Logs de ejecución en tiempo real
    active/             # Sesiones abiertas vinculadas a un plan
    completed/          # Histórico de sesiones cerradas
skills/                 # Capacidades modulares (instrucciones + scripts)
```

---

## 🛠️ Creación de Entidades

Para crear cualquier entidad, se **debe** utilizar el archivo `_template.md` correspondiente ubicado en `clade_skeleton/`.

### 1. Agents (Agentes)
**Template**: `clade_skeleton/agents/agents_template.md`
**Ubicación**: `.claude/agents/{agent-name}.md`

**Cometido**: Definir la personalidad, herramientas y reglas específicas para un rol de IA.
**Posibles Archivos (basado en `claude_examples`):**
- `architect.md`: Especialista en diseño hexagonal y SOLID.
- `developer.md`: Implementación de lógica de negocio y adaptadores.
- `tester.md`: Experto en TDD, unit y e2e testing.
- `reviewer.md`: Validador de calidad de código y estándares arquitectónicos.
- `debugger.md`: Diagnóstico y resolución de bugs complejos.
- `orchestrator-agent.md`: Coordinador maestro del flujo entre agentes.
- `planning-agent.md`: Especialista en desglosar requisitos en tareas técnicas.
- `writer.md`: Encargado de documentación técnica y CHANGELOGs.

### 2. Commands (Comandos)
**Template**: `clade_skeleton/commands/commands_template.md`
**Ubicación**: `.claude/commands/{command-name}.md`

**Cometido**: Automatizar flujos complejos que coordinan múltiples agentes y habilidades.
**Posibles Archivos (basado en `claude_examples`):**
- `bug-fix.md`: Ciclo completo de corrección con reproducción obligatoria.
- `new-feature.md`: Flujo desde el diseño de arquitectura hasta la implementación.
- `refactor.md`: Mejora de estructura sin cambiar comportamiento.
- `pre-flight.md`: Validaciones finales (build, tipos, tests) antes de finalizar.
- `plan-manage.md`: Gestión y actualización de los estados de los planes.
- `orchestrator.md`: El punto de entrada para coordinar otros comandos.

### 3. Skills (Habilidades)
**Template**: `clade_skeleton/skills/skills_template.md`
**Ubicación**: `skills/{skill-name}/SKILL.md` (empaquetado con scripts en `scripts/`)

**Cometido**: Proporcionar conocimientos profundos o herramientas para una tecnología específica.
**Posibles Archivos (basado en `claude_examples`):**
- `node-js.md` / `typescript.md`: Estándares del lenguaje y runtime.
- `express-js.md` / `api-rest.md`: Desarrollo de APIs y middlewares.
- `typeorm.md` / `mysql.md` / `sqlite.md`: Gestión de persistencia y bases de datos.
- `system-architect.md`: Patrones de diseño y arquitectura hexagonal.
- `domain-expert.md` / `usecase-developer.md`: Lógica de negocio y casos de uso.
- `vitest.md` / `supertest.md` / `qa-engineer.md`: Herramientas de testing y QA.
- `zod.md` / `infra-specialist.md`: Validaciones e infraestructura técnica.

### 4. Plans (Planes)
**Template**: `clade_skeleton/plans/plans_template.md`
**Ubicación**: `.claude/plans/active/{task-name}.md`

**Cometido**: Documentar el "Qué", "Por qué" y "Cómo" de una tarea antes de ejecutarla.
- **active/**: Contiene los planes que están siendo ejecutados actualmente.
- **completed/**: Archivo histórico para auditoría y retrospectiva.

### 5. Sessions (Sesiones)
**Template**: `clade_skeleton/sessions/active/sessions_template.md`
**Ubicación**: `.claude/sessions/active/{session-name}.md`

**Cometido**: Registro en tiempo real de las acciones, cambios y métricas de la ejecución de un plan.
- **active/**: Estado actual del trabajo, archivos modificados y bloqueos.
- **completed/**: Registro final de lo que se logró en esa sesión específica.

---

## 📜 Reglas de Oro para Agentes

1. **Uso de Plantillas**: Nunca crear una entidad desde cero; usar siempre los templates de `clade_skeleton/`.
2. **Uso de la Guía**: Todo desarrollo debe alinearse con la `.claude/development_guide.md` (ver `claude_examples`).
3. **Flujo TDD**: La creación de código nuevo o corrección de bugs requiere tests que validen el comportamiento (Fase Roja -> Verde -> Refactor).
4. **Persistencia de Contexto**: Es obligatorio actualizar la **Sesión** tras cada cambio significativo para que otros agentes (o el mismo tras un reinicio) puedan continuar el trabajo sin pérdida de información.
5. **Aprobación de Planes**: Los planes en `plans/active/` deben tener la marca de aprobación del usuario antes de que un agente `developer` comience a escribir código.

---

## 🚀 Requisitos de Scripts (Skills)

Para asegurar la compatibilidad y eficiencia:
- Usar `#!/bin/bash` y `set -e`.
- Salida de estado (logs) a `stderr`.
- Salida de datos (resultados) en **JSON** a `stdout`.
- Referenciar rutas absolutas de scripts como `/mnt/skills/user/{skill-name}/scripts/{script}.sh`.

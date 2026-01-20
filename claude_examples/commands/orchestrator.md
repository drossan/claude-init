---
name: orchestrator
description: Punto de entrada global que analiza la solicitud del usuario y decide qué flujo de trabajo (comando) iniciar para Griddo API. Coordina la transición entre agentes.
usage: "orchestrator [solicitud-del-usuario]"
---

# Comando: Orquestador Global

Este es el comando principal que actúa como cerebro del sistema para Griddo API. Analiza la intención del usuario y delega la ejecución en los comandos especializados, asegurando que se respete la arquitectura hexagonal.

## Lógica de Decisión

El `orchestrator-agent` analizará la solicitud del usuario siguiendo este árbol de decisión:

1.  **¿Es una funcionalidad técnica nueva?** -> Ejecutar `/new-feature`.
2.  **¿Es un error o comportamiento inesperado?** -> Ejecutar `/bug-fix`.
3.  **¿Es una mejora de código sin cambio funcional?** -> Ejecutar `/refactor`.
4.  **¿Faltan tests o hay baja cobertura?** -> Ejecutar `/improve-tests`.
5.  **¿Es una gestión de plan existente?** -> Ejecutar `/plan-manage`.
6.  **¿Es una validación final antes de subir?** -> Ejecutar `/pre-flight`.

## Flujo de Trabajo

### 1. Triaje (Orchestrator Agent)
- Analizar la solicitud.
- Identificar el comando o la combinación de comandos necesarios.
- Informar al usuario sobre el flujo que se va a seguir.
- **Agente**: `orchestrator-agent`

### 2. Delegación
- Invocar el comando especializado correspondiente.
- Asegurar que el contexto y las skills necesarias se pasen correctamente.

### 3. Supervisión
- Monitorizar que cada fase del comando delegado se complete con éxito.
- Intervenir si hay bloqueos que requieran una re-evaluación del flujo.

## Reglas de Oro
- **No adivinar**: Si la solicitud es ambigua, preguntar al usuario antes de decidir el flujo.
- **Enfoque API**: No buscar definiciones de UX/UI; centrarse en endpoints, lógica y datos.
- **Consistencia**: Mantener el archivo maestro de estado en `.claude/sessions/active/global_status.md`.
- **Skills**: Asegurar que los agentes invocados reciban las skills tecnológicas (typescript, node-js, express-js, typeorm, etc.) y técnicas (domain-expert, usecase-developer, etc.) adecuadas.

---

¿En qué puedo ayudarte hoy con Griddo API? Analizaré tu solicitud y activaré el equipo adecuado.

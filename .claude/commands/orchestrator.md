---
name: orchestrator
description: Orquestador maestro que analiza tu solicitud y activa el flujo de trabajo adecuado. Úsalo cuando no estés seguro de qué comando necesitas.
usage: "orchestrator <tu-solicitud>"
---

# Comando: Orquestador (Master Orchestrator)

Este es el comando maestro que analiza tu solicitud y activa el flujo de trabajo adecuado.

## ¿Cuándo Usar?

Usa el orquestador cuando:
- No estés seguro de qué comando necesitas
- Tu solicitud involucre múltiples tipos de trabajo
- Necesites guía sobre el próximo paso

## Flujo de Decisión

El orquestador analiza tu solicitud y decide:

```
¿Es una nueva funcionalidad?
  → Sí: /new-feature

¿Es un error que necesita corrección?
  → Sí: /bug-fix

¿Es mejorar código existente sin cambiar funcionalidad?
  → Sí: /refactor

¿Es mejorar tests?
  → Sí: /improve-tests

¿Es una tarea completada que necesita verificación?
  → Sí: /pre-flight

¿Necesitas ver el estado de un plan?
  → Sí: /plan-manage status
```

## Ejemplos

```bash
# Nueva funcionalidad
claude-init orchestrator "Quiero añadir soporte para proyectos Python"

# Corregir error
claude-init orchestrator "El detector de proyectos no funciona en monorepos"

# Refactorizar
claude-init orchestrator "El código de config tiene demasiadas responsabilidades"

# Verificar
claude-init orchestrator "He terminado de implementar el cliente de Claude"

# Estado
claude-init orchestrator "¿Qué falta para terminar el plan actual?"
```

## Salida

El orquestador te dirá:
1. Qué comando es más apropiado
2. Qué agentes estarán involucrados
3. Qué esperar del proceso

---

¿En qué puedo ayudarte?

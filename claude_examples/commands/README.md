# Guía de Uso de Comandos (Griddo API Automation)

Esta guía explica cuándo utilizar cada comando del sistema de automatización de Griddo API, adaptado a su arquitectura hexagonal.

## 🚀 Comando Maestro: `/orchestrator`
**Úsalo cuando**: No estés seguro de qué comando necesitas.
- Es el punto de entrada recomendado para cualquier tarea.
- Analiza tu solicitud y activa el flujo de trabajo adecuado.

---

## 🏗️ Fase de Construcción: `/new-feature`
**Úsalo cuando**:
- Implementes una funcionalidad nueva en la API.
- Tengas un plan técnico aprobado y quieras empezar a programar siguiendo TDD.
- *Nota: Se enfoca en las capas de Domain, Application e Infrastructure.*

## 🐞 Fase de Mantenimiento: `/bug-fix`
**Úsalo cuando**:
- Tengas un error reportado o un test fallido.
- Quieras investigar la causa raíz de un comportamiento inesperado.
- *Nota: Requiere siempre la creación de un test de reproducción.*

## 🧹 Fase de Mejora: `/refactor`
**Úsalo cuando**:
- Quieras limpiar código, mejorar la legibilidad o aplicar mejores patrones (SOLID).
- No quieras cambiar la funcionalidad externa (la API pública debe ser idéntica).
- *Nota: Se apoya fuertemente en la suite de tests existente.*

## 🧪 Fase de Calidad: `/improve-tests`
**Úsalo cuando**:
- La cobertura de un módulo sea baja.
- Quieras añadir casos de borde a funcionalidades existentes.
- *Nota: Está prohibido modificar código de producción con este comando.*

## 🏁 Fase de Cierre: `/pre-flight`
**Úsalo cuando**:
- Hayas terminado una tarea y quieras asegurarte de que todo es correcto.
- Quieras verificar que no hay "breaking changes" en los contratos de la API.
- Necesites confirmar que el build y los tests pasan satisfactoriamente.

## 📂 Fase de Gestión: `/plan-manage`
**Úsalo cuando**:
- Quieras ver el progreso de una tarea (`status`).
- Quieras retomar un trabajo pausado (`resume`).
- Quieras dar por finalizada una tarea y archivarla (`finish`).

---

### Ejemplo de Flujo Completo:
1. `/orchestrator "Añadir gestión de etiquetas a los sitios"`
2. `/new-feature "Gestión de etiquetas"` (tras aprobación del plan)
3. `/pre-flight`
4. `/plan-manage finish "tags-site"`

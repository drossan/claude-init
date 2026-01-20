# Comandos del claude-init CLI

Esta guía explica cuándo utilizar cada comando del sistema de automatización del claude-init CLI.

## 🚀 Comando Maestro: `/orchestrator`

**Úsalo cuando**: No estés seguro de qué comando necesitas.
- Es el punto de entrada recomendado para cualquier tarea.
- Analiza tu solicitud y activa el flujo de trabajo adecuado.

---

## 🏗️ Desarrollo: `/init`

**Úsalo cuando**:
- Inicialices un nuevo componente del CLI.
- Necesites crear la estructura de paquetes y tests.
- *Nota: Sigue TDD y las convenciones de Go.*

## ✨ Nueva Funcionalidad: `/new-feature`

**Úsalo cuando**:
- Implementes una funcionalidad nueva en el CLI.
- Tengas un plan técnico aprobado.
- *Nota: Requiere un plan aprobado antes de implementar.*

## 🐛 Corrección de Errores: `/bug-fix`

**Úsalo cuando**:
- Tengas un error reportado o un test fallido.
- Quieras investigar la causa raíz.
- *Nota: Requiere crear un test de reproducción.*

## 🧹 Refactorización: `/refactor`

**Úsalo cuando**:
- Quieras mejorar la calidad del código existente.
- No quieras cambiar la funcionalidad externa.
- *Nota: Los tests existentes protegen la funcionalidad.*

## 🧪 Mejora de Tests: `/improve-tests`

**Úsalo cuando**:
- La cobertura de un paquete sea baja.
- Quieras añadir casos de prueba.
- *Nota: Prohibido modificar código de producción.*

## 🏁 Verificación Final: `/pre-flight`

**Úsalo cuando**:
- Hayas terminado una tarea y quieras verificar que todo está correcto.
- Quieras hacer commit o release.
- *Nota: Verifica build, linters, tests y documentación.*

## 📂 Gestión de Planes: `/plan-manage`

**Úsalo cuando**:
- Quieras ver el progreso de una tarea (`status`).
- Quieras retomar un trabajo pausado (`resume`).
- Quieras dar por finalizada una tarea (`finish`).

---

### Ejemplo de Flujo Completo:

1. `/orchestrator "Quiero añadir soporte para la API de OpenAI"`
2. `/new-feature "Soporte para OpenAI"` (tras aprobación del plan)
3. `/pre-flight` (verificación final)
4. `/plan-manage finish openai-support`

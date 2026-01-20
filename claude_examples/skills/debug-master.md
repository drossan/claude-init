# Skill: Maestro de Depuración (debug-master)

## Propósito
Identificar, aislar y resolver la causa raíz de errores, fallos de tests y comportamientos inesperados de manera eficiente y estructurada.

## Responsabilidades
- Realizar análisis de causa raíz (Root Cause Analysis) rastreando el flujo de datos a través de las capas de la arquitectura hexagonal.
- Analizar logs del sistema, trazas de error y el historial de cambios (`git diff`) para localizar el origen del fallo.
- Crear casos de prueba mínimos y deterministas para reproducir bugs antes de intentar corregirlos (metodología "Reproducir antes de Arreglar").
- Utilizar técnicas de depuración avanzada (logging contextual, inspección de estado, validación de mocks).
- Asegurar que las correcciones aplicadas no solo traten los síntomas, sino que resuelvan el problema estructural de fondo.
- Validar que el manejo de excepciones sea coherente y utilice las clases de error de dominio correspondientes.
- Documentar los hallazgos para prevenir la reaparición de problemas similares.

## Skills Tecnológicas Clave
- **typescript**, **node-js**, **express-js**, **typeorm**, **vitest**: Para rastrear y corregir errores en cualquier punto de la arquitectura.

# Plan: Completar SiteFactorySeeder - Generación de Sites Completos

## Descripción
Implementar la generación completa de sites con page data basada en schemas. Se excluye imágenes/archivos por ahora (requiere DAM).

**Aprobado: [x]**
**Fecha de creación:** 2025-01-14
**Agente:** planning-agent

---

## Resumen Ejecutivo

El objetivo es completar el `SiteFactorySeeder` para que genere sites completos basados en schemas. Actualmente el seeder genera:
- ✅ Sites, Languages, DataPacks
- ✅ Taxonomies (categorías)
- ✅ SimpleData (structured data local/global)
- ⚠️ Pages (solo básicas, sin contenido real)
- ❌ Page Data (tipos fromPage: NEWS, EVENT, PROGRAM, etc.)

Este plan se enfoca en implementar **Page Data** para tipos `fromPage: true` del schema.

---

## Estructura del Plan

### Fase 1: Análisis y Diseño 🔍
- [ ] Analizar estructura de schemas para page data
- [ ] Identificar tipos fromPage en el schema
- [ ] Diseñar estrategia de generación de page data
- [ ] Definir interfaz del helper de page data

### Fase 2: Implementación del Helper de Page Data 💻
- [ ] Crear `generatePageData.ts` helper
- [ ] Implementar lógica de generación de campos de schema
- [ ] Implementar asociación con templates
- [ ] Implementar filtrado por dataPacks

### Fase 3: Integración con SiteFactorySeeder 🔗
- [ ] Integrar helper en el seeder principal
- [ ] Manejar datos locales vs globales
- [ ] Implementar limpieza de page data existente

### Fase 4: Testing y Validación ✅
- [ ] Crear tests de integración
- [ ] Validar generación de sites completos
- [ ] Documentar trabajo futuro (imágenes/archivos)

---

## Detalle de Subtareas

### Subtarea 1.1: Análisis de Tipos fromPage

**Tipos identificados en el schema:**
```
QA_GLOBAL_PAGE_DATA, QA_LOCAL_PAGE_DATA, QA_OTHER_GLOBAL_PAGE_DATA,
EVENT, INS_NEWS, NEWS, PRESS_RELEASES, PROGRAM
```

**Cada tipo contiene:**
- `schema.fields[]` - Array de campos con sus definiciones
- `schema.templates[]` - Templates asociados (ej: "NewsDetail")
- `dataPacks[]` - DataPacks que incluyen este tipo
- `local: boolean` - true para local, false para global
- `fromPage: true` - Indica que se genera desde una página

**Ejemplo de schema (NEWS):**
```json
{
  "dataPacks": ["NEWS"],
  "title": "News",
  "local": true,
  "fromPage": true,
  "translate": true,
  "taxonomy": false,
  "schema": {
    "templates": ["NewsDetail"],
    "fields": [
      {
        "type": "TextField",
        "key": "title",
        "mandatory": true
      },
      {
        "type": "DateField",
        "key": "newsDate",
        "indexable": true
      },
      {
        "type": "Wysiwyg",
        "key": "abstract"
      },
      {
        "type": "ImageField",
        "key": "image",
        "mandatory": true
      },
      {
        "type": "MultiCheckSelect",
        "key": "centers",
        "source": "CENTER"
      }
    ]
  }
}
```

---

### Subtarea 2.1: Crear Helper `generatePageData.ts`

**Ubicación:** `src/infrastructure/db/factories/helpers/generatePageData.ts`

**Interfaz propuesta:**

```typescript
interface GeneratePageDataParams {
  dataSource: DataSource;
  siteId: number;
  languageId: number;
  structuredDataType: string;  // Tipo fromPage (NEWS, EVENT, etc.)
  languages: [string, any][];
  pagesToGenerate?: number;
}

interface PageDataResult {
  pageData: Page[];
  count: number;
}

export const generatePageData = async ({
  dataSource,
  siteId,
  languageId,
  structuredDataType,
  languages,
  pagesToGenerate = 1,
}: GeneratePageDataParams): Promise<PageDataResult> => {
  // Implementación...
};
```

**Funcionalidades requeridas:**

1. **Obtener schema del tipo desde schemas.json**
2. **Verificar que el tipo esté en los DataPacks activos del site**
3. **Obtener los templates del schema**
4. **Generar contenido basado en schema.fields[]**
   - Usar `generateContent()` helper existente
   - EXCLUIR campos ImageField y FileField (trabajo futuro)
5. **Crear Páginas con el contenido generado**
   - Asignar template del schema
   - Asignar título y otros campos básicos
   - Guardar el contenido en `page.content` o similar
6. **Retornar las páginas creadas**

---

### Subtarea 2.2: Helper `getPageDataItems.ts`

**Ubicación:** `src/infrastructure/db/factories/helpers/getPageDataItems.ts`

Similar a `getSimpleDataItems.ts` y `getTaxonomyItems.ts`:

```typescript
type ConfigItem = {
  dataPacks?: string[];
  title: string;
  local?: boolean;
  fromPage?: boolean;
  schema?: any;
  [key: string]: any;
};

export type ConfigMap = Record<string, ConfigItem>;

export function getPageDataItems(
  config: ConfigMap,
  local = false,
): ConfigMap {
  return Object.entries(config)
    .filter(([, value]) => {
      return local
        ? value.schema &&
          Object.keys(value.schema).length > 0 &&
          value.fromPage === true &&
          value.local === true
        : value.schema &&
          Object.keys(value.schema).length > 0 &&
          value.fromPage === true &&
          value.local === false;
    })
    .reduce<ConfigMap>((acc, [key, value]) => {
      acc[key] = value;
      return acc;
    }, {});
}
```

---

### Subtarea 3.1: Modificar SiteFactorySeeder

**Archivo:** `src/infrastructure/db/factories/seeders/SiteFactorySeeder.ts`

**Cambios requeridos:**

1. **Importar nuevos helpers:**
```typescript
import { getPageDataItems } from "@infrastructure/db/factories/helpers/getPageDataItems";
import { generatePageData } from "@infrastructure/db/factories/helpers/generatePageData";
```

2. **Agregar lógica de Page Data en el método `seed()`:**

```typescript
async seed(dataSource: DataSource): Promise<void> {
  // ... código existente ...

  for (let index = 0; index < NUM_SITES; index++) {
    const siteId = index + 1;
    // ... código existente ...

    // ✅ NUEVO: Generar Page Data locales para este Site
    const schemaStructuredData = schema.contentTypes.structuredData;
    const localPageDataItems = getPageDataItems(schemaStructuredData, true);

    for (const [pageDataKey, pageDataConfig] of Object.entries(localPageDataItems)) {
      // Verificar si este tipo está en los DataPacks activos del site
      const isInActiveDataPacks = theme.elements?.include?.datapacks?.some(
        (dp: string) => pageDataConfig.dataPacks?.includes(dp)
      );

      if (isInActiveDataPacks) {
        await generatePageData({
          dataSource,
          siteId,
          languageId: 1, // Por ahora, luego iterar por idiomas
          structuredDataType: pageDataKey,
          languages,
          pagesToGenerate: getRandomInt(5, 15),
        });
      }
    }
  }

  // ✅ NUEVO: Generar Page Data GLOBALES
  const globalPageDataItems = getPageDataItems(schemaStructuredData, false);

  for (const [pageDataKey, pageDataConfig] of Object.entries(globalPageDataItems)) {
    await generatePageData({
      dataSource,
      siteId: 0, // 0 para global
      languageId: 1,
      structuredDataType: pageDataKey,
      languages,
      pagesToGenerate: getRandomInt(10, 20),
    });
  }
}
```

---

### Subtarea 4.1: Testing y Validación

**Tests a crear:**

1. **Test unitario del helper:**
   `src/infrastructure/db/factories/__tests__/generatePageData.test.ts`

2. **Test de integración del seeder:**
   `src/infrastructure/db/factories/__tests__/SiteFactorySeeder.pageData.test.ts`

3. **Validación manual:**
   ```bash
   # Ejecutar seeder
   yarn ts-node scripts/factories/qa.ts

   # Verificar en base de datos
   # - Tabla "pages" debe tener registros con content basado en schemas
   # - Los tipos NEWS, EVENT, etc. deben tener páginas
   ```

---

## Trabajo Futuro (NO Implementar en este Plan) 🚫

### Imágenes (ImagesInuse)
- Requiere integración con DAM (Digital Asset Management)
- Subir imágenes y obtener URLs
- Gestionar metadata de imágenes
- Crear tabla ImagesInuse con referencias

### Archivos/Documents
- Requiere integración con DAM
- Subir documentos (PDF, DOC, etc.)
- Gestionar metadata de documentos

**Nota:** Estas features se abordarán en un plan separado cuando el DAM esté disponible.

---

## Checklist de Finalización

- [x] Plan aprobado por usuario (marcar `Aprobado: [x]`)
- [x] Helper `getPageDataItems.ts` creado
- [x] Helper `generatePageData.ts` creado
- [x] `SiteFactorySeeder.ts` modificado
- [x] Build validado sin errores
- [x] Tests corregidos y pasando
- [x] Pre-flight completado
- [x] Validación de breaking changes (sin cambios)
- [ ] Validación manual en QA (requiere configuración de DB - pendiente)
- [x] Plan movido a `.claude/plans/completed/`

**Fecha de finalización:** 2025-01-14

---

## Archivos del Plan

- **Plan:** `.claude/plans/completar-sitefactory-seeder-page-data.md`
- **Sesión:** `.claude/sessions/active/sitefactory-seeder-page-data.md` (se creará al aprobar)

---

## Notas

- Mantener rangos aleatorios (`getRandomInt()`) para cantidad de datos
- NO implementar generación de imágenes/archivos en este plan
- Seguir arquitectura hexagonal del proyecto
- Usar helpers existentes cuando sea posible (`generateContent()`, `createFakePage()`, etc.)

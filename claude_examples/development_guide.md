# Guía de Desarrollo y Contrato de Arquitectura Hexagonal - Griddo API

Este documento establece los estándares, convenciones y el "contrato" de desarrollo para el proyecto Griddo API. Su objetivo es garantizar la consistencia, escalabilidad y mantenibilidad del código a medida que avanzamos en el refactor hacia la Arquitectura Hexagonal.

---

## 📚 Tabla de Contenidos

1. [Estructura de Capas y Directorios](#️-estructura-de-capas-y-directorios)
2. [Estructura de Tests](#-estructura-de-tests)
3. [Monorepo y Gestión de Paquetes](#-monorepo-y-gestión-de-paquetes)
4. [Flujo de Creación de un Endpoint](#-flujo-de-creación-de-un-endpoint-paso-a-paso)
5. [Tipos, Errores y Validaciones](#-tipos-errores-y-validaciones)
6. [Testing Framework y Herramientas](#-testing-framework-y-herramientas)
7. [Patrones de Código Reutilizables](#-patrones-de-código-reutilizables)
8. [Troubleshooting Común](#-troubleshooting-común)
9. [Servicios de Comunicación](#-servicios-de-comunicación)
   - [CommandBus](#commandbus-comunicación-síncrona)
   - [EventBus](#eventbus-comunicación-asíncrona)
   - [Comparación](#comparación-de-servicios-de-comunicación)
10. [Sistema de Colas y Background Jobs](#-sistema-de-colas-y-background-jobs)
11. [AI Services](#-ai-services)
12. [Cache Services](#-cache-services)
13. [Resilience Services](#️-resilience-services)
14. [Sistema de Permisos](#-sistema-de-permisos)
15. [Lifecycle Services](#-lifecycle-services)
16. [Referencias Rápidas](#-referencias-rápidas)
17. [Checklist de Calidad para PRs](#-checklist-de-calidad-para-prs)

---

## 🏛️ Estructura de Capas y Directorios

El código se organiza en tres capas principales, con una dirección de dependencia hacia el núcleo (Domain).

### 1. Capa de Dominio (`src/domain/`)
Es el corazón del sistema, libre de dependencias externas y frameworks.
- **`entities/`**: Objetos de negocio con lógica y validación interna. (Ej: `AnalyticsScriptEntity`).
  - Deben ser clases que validen su propio estado en el constructor.
  - No deben depender de librerías de infraestructura (como TypeORM).
- **`dto/`**: Definiciones de tipos para transferencia de datos entre capas. Siempre `Readonly`.
- **`interfaces/`**: Contratos (Ports) que definen el comportamiento de los adaptadores.
  - **`repositories/`**: Interfaces que debe implementar la infraestructura de persistencia.
  - **`services/`**: Interfaces para servicios externos (Email, Bus de eventos, etc.).
- **`errors/`**: Excepciones específicas del dominio. Deben ser semánticas (ej: `SiteNotFoundError`).
- **`services/`**: Lógica de dominio que no pertenece a una única entidad y no requiere infraestructura.

### 2. Capa de Aplicación (`src/application/`)
Orquesta el flujo de trabajo y ejecuta los casos de uso.
- **`usecases/`**: Clases con una única responsabilidad y un método `execute()`.
  - Reciben dependencias por constructor (Inyección de Dependencias).
  - Orquestan el dominio y los puertos.
  - No conocen detalles de HTTP o Base de Datos.
- **`jobs/`**: Tareas en segundo plano o procesos programados.
- **`services/`**: Servicios de aplicación que coordinan múltiples casos de uso.
- **`mappers/`**: Transformadores de entidades a DTOs de salida.

### 3. Capa de Infraestructura (`src/infrastructure/`)
Implementaciones técnicas y detalles de framework.
- **`controllers/`**: Manejadores de peticiones HTTP.
  - Su única responsabilidad es extraer datos, llamar al caso de uso y responder usando `HttpResponseHandler`.
  - Ejemplo: `/Users/danielrossellosanchez/Documents/Desarrollo/griddo/packages/griddo-api/src/infrastructure/controllers/Site/site.controller.ts`
- **`db/`**: Detalle de persistencia.
  - **`entities/`**: Entidades ORM de TypeORM.
  - **`repositories/`**: Implementaciones reales de las interfaces del dominio.
    - Ejemplo: `SiteRepository` en `/Users/danielrossellosanchez/Documents/Desarrollo/griddo/packages/griddo-api/src/infrastructure/db/repositories/Site/SiteRepository.ts`
  - **`migrations/`**: Scripts de evolución de la base de datos.
  - **`factories/`**: Generadores de datos para testing y seeders.
    - Ejemplo: `createFakeEntitySite` en `/Users/danielrossellosanchez/Documents/Desarrollo/griddo/packages/griddo-api/src/infrastructure/db/factories/definitions/site.factory.ts`
- **`routes/`**: Definición de endpoints y aplicación de middlewares.
  - Ejemplo: `/Users/danielrossellosanchez/Documents/Desarrollo/griddo/packages/griddo-api/src/infrastructure/routes/Site/site.routes.ts`
- **`adapters/`**: Integraciones con servicios externos (Swagger, Bus de eventos, etc.).
  - **`swagger/`**: Configuración y rutas documentadas de OpenAPI.
    - Ejemplo: `/Users/danielrossellosanchez/Documents/Desarrollo/griddo/packages/griddo-api/src/infrastructure/adapters/swagger/routes/page/page.ts`
  - **`type-orm/`**: Adaptador de TypeORM.
  - **`express/`**: Configuración del servidor Express.
- **`dto/`**: Esquemas de **Zod** para validación de entrada/salida y transformación.
  - Ejemplo: `PostSiteDtoSchema` en `/Users/danielrossellosanchez/Documents/Desarrollo/griddo/packages/griddo-api/src/infrastructure/dto/Site/postDTO.ts`
- **`services/`**: Servicios técnicos (cache, configuración, contexto, etc.).
- **`utils/`**: Utilidades de infraestructura (PathUtils, etc.).

---

## 🧪 Estructura de Tests

El proyecto utiliza **Vitest** como framework de testing. La estructura de tests refleja la arquitectura hexagonal del proyecto.

### Organización de Directorios de Tests

```
__tests__/
├── unit/                          # Tests unitarios aislados
│   ├── application/              # Tests de casos de uso
│   │   ├── usecases/             # Tests por caso de uso
│   │   │   ├── Sites/           # create.usescase.test.ts, update.test.ts
│   │   │   ├── Cache/           # clear-distributor-cache.usecase.test.ts
│   │   │   ├── Images/          # FindImagesInFolderUsecase.test.ts
│   │   │   └── ...
│   │   ├── jobs/                # Tests de jobs en segundo plano
│   │   ├── services/            # Tests de servicios de aplicación
│   │   └── utils/               # Tests de utilidades
│   ├── domain/                   # Tests de entidades de dominio
│   │   └── entities/            # Tests de validación y lógica de entidades
│   ├── infrastructure/           # Tests de adapters técnicos
│   │   ├── adapters/            # Tests de adaptadores externos
│   │   ├── db/                  # Tests de repositorios
│   │   │   └── repository/      # Tests de implementaciones de repositorios
│   │   ├── services/            # Tests de servicios técnicos
│   │   └── utils/               # Tests de utilidades de infra
│   └── middleware/               # Tests de middlewares Express
└── integration/                  # Tests de integración
    ├── db/                       # Tests de repos con DB real
    │   └── repository/          # SiteActivity/index.spec.ts
    ├── infrastructure/           # Tests de integración de infra
    │   ├── controllers/         # Tests de controladores con HTTP
    │   ├── routes/              # Tests de rutas completas
    │   └── services/            # Tests de servicios con dependencias
    └── adapters/                # Tests de adaptadores externos
```

### Patrones de Testing

#### 1. Patrón AAA (Arrange-Act-Assert)

Todos los tests deben seguir el patrón AAA:

```typescript
// Ejemplo real: __tests__/unit/application/usecases/Sites/create.usescase.test.ts:29-60
describe("CreateSiteUsesCase", () => {
    let repository: Partial<SiteRepositoryInterface>;
    let eventBus: Partial<EventBusInterface>;
    let usesCase: CreateSiteUsesCase;

    beforeEach(() => {
        // ARRANGE: Preparar el entorno
        repository = {
            create: vi.fn(),
        };
        eventBus = {
            emit: vi.fn().mockResolvedValue(undefined),
        };
        usesCase = new CreateSiteUsesCase(
            repository as SiteRepositoryInterface,
            eventBus as EventBusInterface,
        );
    });

    it("should create a site via repository and return DTO", async () => {
        // ARRANGE: Preparar datos de prueba
        const createdEntity = createFakeEntitySite({
            id: 101,
            name: "My New Site",
            authorId: 7,
            slug: "my-new-site",
        });
        repository.create = vi.fn().mockResolvedValueOnce(createdEntity);

        // ACT: Ejecutar el código a testear
        const result = await usesCase.execute({
            name: "My New Site",
            authorId: 7,
            slug: "my-new-site",
        });

        // ASSERT: Verificar resultados
        expect(repository.create).toHaveBeenCalledTimes(1);
        expect(result).toEqual(createdEntity.toDTO());
        expect(eventBus.emit).toHaveBeenCalled();
    });
});
```

#### 2. Mocks con Vitest

Usa `vi.fn()` y `vi.mock()` para crear dobles de prueba:

```typescript
// Mock de repositorio
repository = {
    create: vi.fn(),
    findById: vi.fn(),
    update: vi.fn(),
};

// Configurar comportamiento del mock
repository.create = vi.fn().mockResolvedValueOnce(fakeEntity);

// Verificar llamadas al mock
expect(repository.create).toHaveBeenCalledWith(expect.any(SiteEntity));
expect(repository.create).toHaveBeenCalledTimes(1);
```

#### 3. Factories para Datos de Prueba

Usa las factories definidas en `src/infrastructure/db/factories/definitions/`:

```typescript
import { createFakeEntitySite } from "@infrastructure/db/factories/definitions/site.factory";

// Crear entidad con datos aleatorios
const site = createFakeEntitySite();

// Crear entidad con overrides
const site = createFakeEntitySite({
    id: 101,
    name: "Custom Name",
    authorId: 7,
});
```

#### 4. Tests de Integración con Base de Datos

Los tests de integración usan una base de datos SQLite en memoria:

```typescript
// Ejemplo real: __tests__/integration/db/repository/Sites/SiteActivity/index.spec.ts:7-28
describe("SiteActivityRepository", () => {
    let repository: SiteActivityRepositoryInterface;
    let db: DBTypeORM;

    beforeEach(async () => {
        db = await DBTypeORM.getInstance();
        repository = new SiteActivityRepository(db);
    });

    it("should insert a new SiteActivity and return its ID", async () => {
        const dto = createFakeSiteActivity();
        const entity = await repository.save(dto);

        expect(entity).toBeDefined();
        expect(entity.userId).toBe(dto.userId);
    });
});
```

### Configuración de Tests

#### Setup Global (`vitest.setup.ts`)

El archivo `vitest.setup.ts` configura:
- Base de datos de prueba (SQLite en memoria)
- Seeders iniciales (usuarios, sites, eventos, etc.)
- Configuración global (CONFIG, environment variables)
- Inicialización de EventBus y CommandBus

#### Configuración de Vitest (`vitest.config.ts`)

```typescript
{
    test: {
        environment: "node",           // Entorno Node.js
        globals: true,                  // APIs globales (describe, it, etc.)
        clearMocks: true,               // Limpiar mocks entre tests
        include: ["**/__tests__/**/*.(test|spec).(ts|js)"],
        setupFiles: ["./vitest.setup.ts"],
        coverage: {
            provider: "v8",             // Proveedor de cobertura
            reporter: ["text", "json", "html"],
        },
    },
}
```

### Ejecutar Tests

```bash
# Todos los tests
yarn test

# Tests en modo watch
yarn test:watch

# Tests con cobertura
yarn test:coverage

# Solo tests de integración
yarn test:integration
```

### Mejores Prácticas de Testing

1. **Tests Unitarios**: Deben ser rápidos y aislados, sin dependencias externas.
2. **Tests de Integración**: Prueban la integración entre capas, usan base de datos real.
3. **Nombre de Tests**: Debe ser descriptivo: `"should [expected behavior] when [condition]"`
4. **Un Test por Caso**: Cada `it` debe probar un único comportamiento.
5. **Setup/Teardown**: Usa `beforeEach`/`afterEach` para limpiar estado entre tests.
6. **Factories**: Siempre usa factories en lugar de crear datos manualmente.

---

## 📦 Monorepo y Gestión de Paquetes

### Estructura del Monorepo

```
packages/
├── griddo-api/           # API principal (este proyecto)
├── griddo-api-public/    # API pública
├── griddo-ax/            # Admin X
├── griddo-cx/            # Customer X
├── griddo-components/    # Componentes compartidos
├── griddo-core/          # Core compartido
├── api-types/            # Tipos compartidos entre frontend y backend
└── eslint-config-griddo-back/ # Configuración de ESLint
```

### Paquete `api-types`

El paquete `@griddo/api-types` contiene los tipos públicos de la API que son compartidos entre:

- **Backend** (`griddo-api`): Para tipado de respuestas y solicitudes
- **Frontend** (AX/CX): Para tipado de llamadas a la API
- **SDK**: Para generar clientes tipados

#### Estructura de `api-types`

```
packages/api-types/src/
├── common/               # Tipos comunes (paginación, exportación)
├── logs/                 # Tipos de logs y actividad
├── schemas/              # Tipos de esquemas de datos
├── site/                 # Tipos de sitios
├── pages/                # Tipos de páginas
├── redirects/            # Tipos de redirecciones
├── dataPack/             # Tipos de data packs
├── queue/                # Tipos de colas
└── index.ts              # Exportaciones públicas
```

#### Regla de Oro

> **Si un tipo/interface es consumido por el frontend (AX/CX) o el SDK, DEBE vivir en `packages/api-types`.**

El proyecto principal importa estos tipos para garantizar la sincronización automática de los contratos de la API.

#### Uso de `api-types`

```typescript
// Importar tipos compartidos
import type { SiteDTO, SitePostParamsDTO, SitePostResponseDTO } from "packages/api-types/src";

// Usar en entidades de dominio
export class SiteEntity {
    toDTO(): SiteDTO {
        return {
            id: this.id?.getValue(),
            name: this.name,
            // ...
        };
    }
}

// Usar en controladores
static async post(
    request: Request<{}, {}, SitePostParamsDTO>,
    response: Response<SitePostResponseDTO>,
): Promise<void> {
    // ...
}
```

#### Publicación de `api-types`

```bash
# Entrar al directorio del paquete
cd packages/api-types

# Construir
npm run build

# Publicar (npm publish se encarga de versionado)
npm publish
```

### Dependencias entre Paquetes

- `griddo-api` NO debe depender de `griddo-ax` o `griddo-cx`
- `griddo-ax` y `griddo-cx` pueden depender de `api-types`
- `griddo-core` y `griddo-components` son compartidos por todos los proyectos frontend

### Gestión de Versiones

- Usa `yarn workspaces` para gestionar dependencias
- Las versiones se sincronizan mediante cambios en `package.json`
- `api-types` tiene versionado semántico independiente

---

## 🚀 Flujo de Creación de un Endpoint (Paso a Paso)

Para añadir una nueva funcionalidad siguiendo TDD, sigue este flujo:

### 1. Definir el Contrato en el Dominio
- Crea el DTO de entrada/salida en `src/domain/dto/`.
- Define la interfaz del repositorio en `src/domain/interfaces/`.

### 2. Crear la Entidad de Dominio
- Implementa la lógica de negocio y validación interna.
- Crea el test unitario de la entidad.

### 3. Implementar el Caso de Uso (TDD)
- Crea el test del caso de uso en `__tests__/unit/application/usecases/`.
- Implementa el Caso de Uso en `src/application/usecases/` inyectando la interfaz del repositorio.

### 4. Implementar la Persistencia
- Crea la entidad de TypeORM en `src/infrastructure/db/entities/`.
- Implementa el repositorio en `src/infrastructure/db/repositories/`.
- Crea el test de integración del repositorio.

### 5. Crear el Controlador y la Ruta
- Define el esquema de validación **Zod** en `src/infrastructure/dto/`.
- Crea el controlador en `src/infrastructure/controllers/` usando `DatabaseConnector.run`.
- Define la ruta en `src/infrastructure/routes/`.

### 6. Documentación Mandatoria en Swagger
Todo endpoint **DEBE** estar documentado.
- Crea el archivo de ruta Swagger en `src/infrastructure/adapters/swagger/routes/`.
- Usa `zod-to-openapi` para vincular los esquemas de Zod con la documentación.
- Registra la nueva ruta en el `index.ts` del adaptador de Swagger.

**Ejemplo real** (`src/infrastructure/adapters/swagger/routes/page/page.ts:21-104`):

```typescript
export const registerPageRoutes = (registry: OpenAPIRegistry) => {
    registry.registerPath({
        method: "get",
        path: "/select/site/{site}/pages/{excludePage}",
        tags: ["Pages"],
        summary: "Select pages for a site",
        security: [{ bearerAuth: [] }],
        parameters: [
            {
                name: "site",
                in: "path",
                required: true,
                schema: { type: "integer", format: "int64" },
                description: "Site ID",
            },
        ],
        responses: {
            200: {
                description: "List of pages",
                content: {
                    "application/json": {
                        schema: PagesSelectResponseSchema,
                    },
                },
            },
        },
    });
};
```

---

## 🛠️ Tipos, Errores y Validaciones

### Gestión de Tipos Compartidos (`api-types`)
- **Regla de Oro**: Si un tipo/interface es consumido por el frontend (AX/CX) o el SDK, **DEBE** vivir en `packages/api-types`.
- El proyecto principal importa estos tipos para garantizar la sincronización automática de los contratos de la API.

### Errores del Dominio

Los errores específicos del dominio extienden de `Error` y proporcionan contexto semántico:

**Ejemplo real** (`src/domain/errors/Site/index.error.ts:1-6`):

```typescript
export class SiteError extends Error {
    constructor(message: string, prefix = "[Site] ") {
        super(`${prefix}${message}`);
        this.name = "SiteError";
    }
}
```

**Uso en entidades** (`src/domain/entities/Site/index.entity.ts:96-104`):

```typescript
private validate() {
    if (!this.name) {
        throw new SiteError("Name is required");
    }

    if (this.name.length > 255) {
        throw new SiteError("Name cannot exceed 255 characters");
    }
}
```

### Value Objects

Los Value Objects encapsulan validaciones de valores primitivos:

**Ejemplo real** (`src/domain/valueObjects/id_number.valueobject.ts:1-25`):

```typescript
export class IdNumber {
    private readonly value: number;

    constructor(id: number) {
        if (!Number.isInteger(id) || id <= 0) {
            throw new Error(
                `ID must be a positive integer greater than 0. Received: ${id}`,
            );
        }
        this.value = id;
    }

    getValue(): number {
        return this.value;
    }

    equals(other: IdNumber): boolean {
        return this.value === other.getValue();
    }
}
```

### Validación con Zod

Zod se utiliza en la capa de infraestructura para validar entrada/salida:

**Ejemplo real** (`src/infrastructure/dto/Site/postDTO.ts:5-22`):

```typescript
export const PostSiteDtoSchema = z.object({
    name: z.string().min(1, "Name is required"),
    domain: z.number().min(1, "Domain is required"),
    defaultLanguage: z.number().min(1, "Language is required"),
    path: z.string().optional().nullable().default(null),
});

// Función de transformación
export const toPostSiteParams = (dto: PostDTO): CreateDTO => {
    return {
        ...dto,
        languageId: dto.defaultLanguage,
    };
};
```

### Gestión de Errores en Controladores

**Ejemplo real** (`src/infrastructure/controllers/Site/site.controller.ts:235-249`):

```typescript
static async post(
    request: Request<{}, {}, SitePostParamsDTO>,
    response: Response<SitePostResponseDTO>,
): Promise<void> {
    await DatabaseConnector.run(response, async (db) => {
        try {
            // ... lógica del controlador
            HttpResponseHandler.ok<SitePostResponseDTO>(response, site);
        } catch (error) {
            if (error instanceof z.ZodError) {
                HttpResponseHandler.error(response, error);
                return;
            }

            if (error instanceof Error) {
                HttpResponseHandler.error(response, error.message);
                return;
            }

            HttpResponseHandler.error(response, "Unknown error");
        }
    });
}
```

---

## 🔧 Testing Framework y Herramientas

### Stack de Testing

- **Framework**: Vitest (alternativa moderna a Jest)
- **Coverage**: V8 (integrado con Vitest)
- **Mocks**: Vitest (`vi.fn()`, `vi.mock()`)
- **Base de datos de测试**: SQLite (better-sqlite3)
- **Factories**: @faker-js/faker para datos aleatorios

### Comandos de Testing

```bash
# Ejecutar todos los tests
yarn test

# Tests en modo watch (desarrollo)
yarn test:watch

# Tests con cobertura de código
yarn test:coverage

# Solo tests de integración
yarn test:integration

# Tests específicos por patrón
yarn test -- Sites
```

### Configuración de Path Aliases

Los tests usan los mismos path aliases que el código fuente (`vitest.config.ts:18-43`):

```typescript
resolve: {
    alias: {
        "@domain": path.resolve(__dirname, "src/domain"),
        "@application": path.resolve(__dirname, "src/application"),
        "@usecases": path.resolve(__dirname, "src/application/usecases"),
        "@infrastructure": path.resolve(__dirname, "src/infrastructure"),
        "@repository": path.resolve(__dirname, "src/infrastructure/db/repositories"),
        "@controllers": path.resolve(__dirname, "src/infrastructure/controllers"),
    },
}
```

---

## 🔄 Patrones de Código Reutilizables

### Patrón Repository

Los repositorios implementan interfaces del dominio y manejan la persistencia:

**Ejemplo completo** (`src/infrastructure/db/repositories/Site/SiteRepository.ts:18-210`):

```typescript
export class SiteRepository implements SiteRepositoryInterface {
    constructor(
        private readonly dataSource: DatabaseAdapter,
        private readonly uuidGenerator?: UuidGeneratorInterface,
    ) {}

    // Método para convertir de TypeORM a Dominio
    private toDomain(site: Site): SiteEntity {
        const siteDTO: SiteDTO = {
            id: site.id,
            authorId: site.authorId,
            name: site.name,
            // ... mapeo de campos
        };
        return new SiteEntity(siteDTO);
    }

    async create(site: SiteEntity): Promise<SiteEntity> {
        const siteDTO = site.toDTO();
        const repository = await this.dataSource.getRepository(Site);

        const siteTypeORM = repository.create({
            authorId: siteDTO.authorId ?? undefined,
            name: siteDTO.name,
            // ... mapeo de campos
        });

        const savedSite = await repository.save(siteTypeORM);
        return this.toDomain(savedSite);
    }

    async findById(id: number): Promise<SiteEntity | null> {
        const repository = await this.dataSource.getRepository(Site);
        const site = await repository.findOne({ where: { id } });
        return site ? this.toDomain(site) : null;
    }
}
```

### Patrón Use Case

Los casos de uso orquestan la lógica de negocio:

**Ejemplo simple** (`src/application/usecases/Site/create.usescase.ts:10-44`):

```typescript
export class CreateSiteUsesCase {
    constructor(
        private readonly repository: SiteRepositoryInterface,
        private readonly eventBus: EventBusInterface,
    ) {}

    async execute({ name, authorId, slug }: SiteDTO) {
        const siteEntity = new SiteEntity({
            name,
            authorId,
            slug,
        });

        const site = await this.repository.create(siteEntity);

        // Emitir evento de forma asíncrona (no bloqueante)
        this.logUserAction(site.toDTO()).catch((error) => {
            console.error(
                `[CreateSiteUsesCase - USER_ACTION] Error: ${error.message}`,
            );
        });

        return site.toDTO();
    }
}
```

**Ejemplo complejo con orquestación** (`src/application/usecases/Site/save.usescase.ts:17-170`):

```typescript
export class SaveSiteUsesCase {
    constructor(
        private readonly repository: SiteRepositoryInterface,
        private readonly pageRepository: PagesRepositoryInterface,
        private readonly languageSiteRepository: SiteLanguageRepositoryInterface,
        private readonly eventBus: EventBusInterface,
        private readonly commandBus: CommandBusInterface,
    ) {}

    async execute(params: SiteDTO & {
        path: string;
        siteId?: number | null;
        languageId?: number | null;
        domainId?: number | null;
    }) {
        let site: SiteDTO | null;

        // 1. Encontrar slug disponible
        const candidateSlug = await this.findAvailableSlug(id, slug);

        // 2. Crear o actualizar según corresponda
        if (!id) {
            const usesCase = new CreateSiteUsesCase(
                this.repository,
                this.eventBus,
            );
            site = await usesCase.execute({ name, authorId, slug: candidateSlug });
        } else {
            const usesCase = new UpdateSiteUsesCase(
                this.repository,
                this.pageRepository,
                this.eventBus,
            );
            site = await usesCase.execute({ /* params */ });
        }

        // 3. Asignar lenguaje al sitio
        if (languageId && domainId) {
            await this.addLanguageToSite({ domainId, path, siteId: site.id!, languageId });
        }

        return site;
    }
}
```

### Patrón Controller

Los controladores manejan HTTP y delegan a casos de uso:

**Ejemplo** (`src/infrastructure/controllers/Site/site.controller.ts:187-251`):

```typescript
export class SiteController {
    static async post(
        request: Request<{}, {}, SitePostParamsDTO>,
        response: Response<SitePostResponseDTO>,
    ): Promise<void> {
        await DatabaseConnector.run(response, async (db) => {
            try {
                // 1. Extraer contexto
                const { authorId, siteId } = SiteController.getContext(request);

                // 2. Validar y transformar parámetros
                const { name, languageId, domain, path } =
                    SiteController.postParamsDTO(request);

                // 3. Validaciones de negocio
                const languageSiteRepository = new SiteLanguageRepository(db);
                await SiteController.checkLanguagePath(
                    languageId, domain, path ?? "", languageSiteRepository,
                );

                // 4. Crear repositorios y caso de uso
                const repositorySite = new SiteRepository(db, new UuidGenerator());
                const pageRepository = new PagesRepository(db, new UuidGenerator());

                const saveSiteUsesCase = new SaveSiteUsesCase(
                    repositorySite,
                    pageRepository,
                    languageSiteRepository,
                    eventBus,
                    commandBus,
                );

                // 5. Ejecutar caso de uso
                const site = await saveSiteUsesCase.execute({
                    name, authorId, siteId, languageId, domainId: domain, path,
                });

                // 6. Responder
                HttpResponseHandler.ok<SitePostResponseDTO>(response, site);
            } catch (error) {
                // Manejo de errores
                if (error instanceof z.ZodError) {
                    HttpResponseHandler.error(response, error);
                    return;
                }
                if (error instanceof Error) {
                    HttpResponseHandler.error(response, error.message);
                    return;
                }
                HttpResponseHandler.error(response, "Unknown error");
            }
        });
    }
}
```

### Patrón DatabaseConnector

El `DatabaseConnector` abstrae la gestión de transacciones y conexiones:

**Ejemplo** (`src/infrastructure/db/DatabaseConnector.ts:17-64`):

```typescript
export class DatabaseConnector {
    static async run(
        res: Response,
        callback: (db: DatabaseAdapter) => Promise<void>,
        { readOnly = false, useORM = true }: ConnectorOptions = {},
    ): Promise<void> {
        // Usar base de datos de test si estamos en entorno de test
        if (ServerState.getEnvironment() === "test") {
            const dataSource = AppTestDataSource;
            DBTypeORM.useCustomDataSource(dataSource);
        }

        try {
            const db = useORM
                ? await DBTypeORM.getInstance(readOnly)
                : DB.getInstance(readOnly);

            await callback(db as unknown as DatabaseAdapter);
        } catch (err) {
            const msg = err instanceof Error ? err.message : "Unknown error";
            outputError(res, msg);
        }
    }

    static runRO(
        res: Response,
        callback: (db: DatabaseAdapter) => Promise<void>,
        params: Omit<ConnectorOptions, "readOnly"> = {},
    ) {
        return this.run(res, callback, { ...params, readOnly: true });
    }
}
```

---

## 🚨 Troubleshooting Común

### Problemas Frecuentes y Soluciones

#### 1. Tests fallan con "Cannot find module"

**Problema**: Los tests no pueden resolver los path aliases.

**Solución**:
- Verifica que `vitest.config.ts` tiene los alias configurados correctamente
- Asegúrate de importar usando los aliases: `@domain`, `@application`, etc.
- Reinicia el servidor de tests

#### 2. Error de tipos al importar desde `api-types`

**Problema**: TypeScript no encuentra los tipos de `packages/api-types`.

**Solución**:
```bash
# Construir el paquete api-types
cd packages/api-types
npm run build

# Volver a la raíz y reinstalar dependencias
cd ../griddo-api
yarn install
```

#### 3. Tests de integración fallan con "Database not initialized"

**Problema**: La base de datos de test no se inicializó correctamente.

**Solución**:
- Verifica que `vitest.setup.ts` se ejecuta antes de los tests
- Revisa que `AppTestDataSource` está correctamente configurado
- Asegúrate de que los seeders se ejecutan sin errores

#### 4. Error "Repository method not implemented"

**Problema**: La interfaz del dominio tiene un método que el repositorio no implementa.

**Solución**:
- Implementa el método faltante en el repositorio de infraestructura
- O, si no es necesario, elimínalo de la interfaz del dominio

#### 5. Zod validation error en producción pero no en tests

**Problema**: Los datos de producción tienen un formato diferente.

**Solución**:
- Revisa el esquema Zod en `src/infrastructure/dto/`
- Añade validaciones más específicas o transforma los datos antes de validar
- Usa `.optional().nullable()` para campos que pueden ser `null` o `undefined`

#### 6. EventBus o CommandBus no funcionan en tests

**Problema**: Los buses no están inicializados o están mockeados incorrectamente.

**Solución**:
```typescript
// Mock correcto de EventBus
eventBus = {
    emit: vi.fn().mockResolvedValue(undefined),
};

// Mock correcto de CommandBus
commandBus = {
    dispatch: vi.fn().mockResolvedValue(undefined),
};
```

#### 7. Migrations fallan al ejecutarse

**Problema**: La migración genera SQL inválido o conflictos con la base de datos.

**Solución**:
- Revisa la sintaxis de TypeORM en la migración
- Ejecuta `yarn migration:show` para ver el estado
- Si es necesario, crea una nueva migración que corrija el problema

#### 8. Performance issues en queries de repositorios

**Problema**: Las consultas son lentas o devuelven demasiados datos.

**Solución**:
- Usa `select` para limitar los campos devueltos
- Añade índices en la base de datos si es necesario
- Usa relaciones con `eager: true` solo cuando sea necesario
- Considera paginación para grandes conjuntos de datos

---

## ✅ Checklist de Calidad para PRs
- [ ] ¿Se han creado/actualizado los tipos en `api-types` si afectan al cliente?
- [ ] ¿El endpoint tiene su esquema Zod de entrada y salida?
- [ ] ¿Está el endpoint registrado y documentado en Swagger?
- [ ] ¿La lógica de negocio está en el Caso de Uso y no en el Controlador?
- [ ] ¿Se inyectan las dependencias por interfaz?
- [ ] ¿Existen tests unitarios para el Caso de Uso?
- [ ] ¿Existen tests de integración para la ruta?
- [ ] ¿Se usa `HttpResponseHandler` y `DatabaseConnector` correctamente?

---

## 🔌 Servicios de Comunicación

El proyecto utiliza tres patrones de comunicación diferentes según el caso de uso:

### CommandBus (Comunicación Síncrona)

**Ubicación**: `src/infrastructure/services/commandBus/`

El **CommandBus** implementa el patrón mediador para comunicación síncrona (solicitud-respuesta). Úsalo cuando necesites una respuesta inmediata.

#### Características

- **Patrón**: Request-Response
- **Sincronía**: Síncrono
- **Respuesta**: Devuelve resultado
- **Manejo de errores**: Propaga errores al emisor

#### Inicialización

```typescript
import { CommandBusService } from "@infrastructure/services/commandBus";

// En el arranque de la aplicación
await CommandBusService.initialize();
```

#### Manejadores Registrados

1. **DependentPagesUpdaterHandler**: Actualiza páginas dependientes
2. **RegisterDistributorStructuredDataHandler**: Registra datos estructurados
3. **FetchDistributorContentHandler**: Obtiene contenido de distribuidores
4. **FetchPageHandler**: Obtiene páginas
5. **ClearSearchesHandler**: Limpia búsquedas en caché
6. **GenerateEmbeddingsHandler**: Genera embeddings
7. **FindStructuredDataHandler**: Busca datos estructurados
8. **SaveSimpleDataSearchesHandler**: Guarda búsquedas

#### Uso en Casos de Uso

```typescript
import type { CommandBusInterface } from "@domain/interfaces/services/CommandBusInterface";

export class SomeUseCase {
    constructor(
        private readonly commandBus: CommandBusInterface,
    ) {}

    async execute(params: Params) {
        // Dispatch síncrono - espera respuesta
        const result = await this.commandBus.dispatch<ResponseType>({
            commandName: "FETCH_PAGE",
            payload: { pageId: params.pageId },
        });

        return result;
    }
}
```

#### Cuándo Usar CommandBus

✅ **Usar cuando:**
- Necesitas una respuesta inmediata
- La operación es crítica para el flujo
- El emisor necesita manejar errores

❌ **No usar cuando:**
- La operación puede tomar mucho tiempo
- Quieres notificar a múltiples componentes
- La respuesta no es necesaria

---

### EventBus (Comunicación Asíncrona)

**Ubicación**: `src/infrastructure/services/eventBus/`

El **EventBus** implementa el patrón publicador/suscriptor para comunicación asíncrona (fire-and-forget). Úsalo para notificaciones y eventos del sistema.

#### Características

- **Patrón**: Publisher-Subscriber
- **Sincronía**: Asíncrono
- **Respuesta**: No devuelve resultado
- **Manejo de errores**: Los handlers manejan sus errores internamente

#### Inicialización

```typescript
import { EventBusService } from "@infrastructure/services/eventBus";

// En el arranque de la aplicación
await EventBusService.initialize();
```

#### Manejadores Registrados

1. **UserActionLogger**: Registra acciones de usuarios para auditoría

> **Nota**: Muchos manejadores fueron migrados al CommandBus. El EventBus ahora se usa principalmente para eventos del sistema.

#### Uso en Casos de Uso

```typescript
import type { EventBusInterface } from "@domain/interfaces/services/EventBusInterface";

export class CreateSiteUsesCase {
    constructor(
        private readonly repository: SiteRepositoryInterface,
        private readonly eventBus: EventBusInterface,
    ) {}

    async execute(params: SiteDTO) {
        const site = await this.repository.create(siteEntity);

        // Emitir evento de forma NO bloqueante
        this.eventBus.emit({
            eventName: "SITE_CREATED",
            payload: { siteId: site.id, name: site.name },
        }).catch((error) => {
            console.error(`[EventBus] Error: ${error.message}`);
        });

        return site.toDTO();
    }
}
```

#### Cuándo Usar EventBus

✅ **Usar cuando:**
- Quieres notificar a múltiples componentes
- La respuesta no es necesaria
- El evento es de auditoría o logging

❌ **No usar cuando:**
- Necesitas una respuesta
- La operación es crítica para el flujo
- El emisor necesita saber si falló

---

### Comparación de Servicios de Comunicación

| Característica | CommandBus | EventBus | Queue System |
|---------------|-----------|----------|--------------|
| Patrón | Request-Response | Publisher-Subscriber | Producer-Consumer |
| Sincronía | Síncrono | Asíncrono inmediato | Asíncrono diferido |
| Respuesta | Devuelve resultado | No devuelve respuesta | No devuelve respuesta |
| Persistencia | No | No | Sí (SQS/RabbitMQ) |
| Reintentos | No | No | Automáticos |
| Programación | Manual | Event-driven | Cron jobs |
| Uso típico | Operaciones inmediatas | Notificaciones | Tareas pesadas/programadas |

---

## 📋 Sistema de Colas y Background Jobs

**Ubicación**: `src/application/jobs/` y `src/infrastructure/services/queue/`

El sistema de colas permite ejecutar tareas en segundo plano de forma asíncrona y programada.

### Componentes Principales

#### 1. QueueWorkerService

Servicio principal que gestiona el sistema de colas.

```typescript
import { QueueWorkerService } from "@infrastructure/services/queue/QueueWorkerService";

// Inicialización
await QueueWorkerService.initialize();
await QueueWorkerService.startProcessing();

// Dispatch de trabajos
QueueWorkerService.getInstance().dispatch("JOB_NAME", payload);
```

#### 2. BackgroundJobService

Motor de procesamiento que implementa:
- Registro de handlers
- Dispatch con prioridades
- Procesamiento con concurrencia controlada
- Reintentos automáticos

### Tipos de Trabajos Registrados

#### Trabajos de Páginas (`pageJobHandlers`)

- **CLEAR_EDITING_PAGES**: Limpia estados de edición caducados
  - Intervalo: 1 minuto
  - Instancia única: Sí

- **EMIT_RELATED_PAGE_UPDATE**: Emite actualizaciones a páginas relacionadas
  - Disparado por eventos
  - Procesamiento bajo demanda

#### Trabajos de Caché (`cacheJobHandlers`)

- **DELETED_EXPIRED_CACHED**: Borra datos expirados
  - Intervalo: 12 horas
  - Instancia única: Sí

#### Trabajos de Búsqueda (`searchJobHandlers`)

- **CLEAR_SEARCHES**: Limpia búsquedas en caché
  - Intervalo: 15 minutos
  - Instancia única: Sí

#### Trabajos de IA (`AIJobHandlers`)

- **LOAD_EMBEDDINGS**: Carga embeddings pendientes
  - Intervalo: 60 minutos
  - Todas las instancias: Sí

### Cómo Crear un Nuevo Job

#### 1. Crear el Handler

```typescript
// src/application/jobs/newCategory/NewJobHandler.ts
import type {
    BackgroundJobHandler,
    JobContext,
} from "@domain/interfaces/services/BackgroundJobInterface";

export class NewJobHandler implements BackgroundJobHandler<NewJobPayload> {
    async handle(payload: NewJobPayload, context: JobContext): Promise<void> {
        const db = await DBTypeORM.getInstance();
        const repository = new SomeRepository(db);

        await repository.performSomeOperation(payload.data);

        console.info(`Job completed: ${context.jobId}`);
    }
}
```

#### 2. Definir el Payload

```typescript
// src/domain/dto/Queue/Handlers/NewJobPayload.ts
export interface NewJobPayload {
    readonly data: string;
    readonly options?: {
        readonly retries?: number;
        readonly priority?: "high" | "normal" | "low";
    };
}
```

#### 3. Registrar en el Enum

```typescript
// src/application/jobs/jobsHandler.ts
export enum jobHandlers {
    NEW_JOB_NAME = "NEW_JOB_NAME",
}
```

#### 4. Exportar en su Categoría

```typescript
// src/application/jobs/newCategory/index.ts
import { NewJobHandler } from "./NewJobHandler";

export const newCategoryJobHandlers = {
    NEW_JOB_NAME: new NewJobHandler(),
};
```

#### 5. Incluir en Handlers Globales

```typescript
// src/application/jobs/index.jobs.ts
import { newCategoryJobHandlers } from "@application/jobs/newCategory";

export const allJobHandlers = {
    ...newCategoryJobHandlers,
};
```

#### 6. Configurar Cron Job (opcional)

```typescript
// common/cron.ts
setCron(
    "Descripción del trabajo",
    () => QueueWorkerService.getInstance().dispatch(jobHandlers.NEW_JOB_NAME, payload),
    { useDatabase: true },
    30, // cada 30 minutos
    true, // solo instancia principal
);
```

### Dispatch Manual

```typescript
// En controladores, casos de uso, etc.
import { QueueWorkerService } from "@infrastructure/services/queue/QueueWorkerService";
import { jobHandlers } from "@application/jobs/jobsHandler";

// Dispatch inmediato
await QueueWorkerService.getInstance().dispatch(jobHandlers.NEW_JOB_NAME, {
    data: "payload-data",
});

// Dispatch con opciones
await QueueWorkerService.getInstance().dispatch(
    jobHandlers.NEW_JOB_NAME,
    { data: "payload-data" },
    {
        priority: "high",
        delay: 5000, // 5 segundos
        attempts: 3, // máximo 3 reintentos
    },
);
```

### Cuándo Usar el Sistema de Colas

✅ **Usar cuando:**
- La tarea puede tomar >30 segundos
- Se requiere procesamiento programado/recurrente
- La tarea puede fallar y necesita reintentos
- Quieres desacoplar la ejecución del request HTTP

❌ **No usar cuando:**
- Se necesita respuesta inmediata
- La operación es simple y rápida (<1 segundo)
- Se requiere procesamiento síncrono

---

## 🤖 AI Services

**Ubicación**: `src/infrastructure/services/AI/`

Servicios para integración con proveedores de inteligencia artificial.

### AIServiceFactory

Fábrica que proporciona instancias de servicios de IA según el adaptador especificado.

```typescript
import { AIAdapters } from "@infrastructure/adapters/ai/AIAdapters";
import { AIServiceFactory } from "@infrastructure/services/AI/AIServiceFactory";

// Obtener servicio de OpenAI
const aiService = AIServiceFactory.getAdapter(AIAdapters.OPENAI);

// Utilizar el servicio
const response = await aiService.generateText(prompt);
```

### Adaptadores Soportados

- **OpenAI** (`AIAdapters.OPENAI`): Generación de texto, embeddings

### Cómo Añadir un Nuevo Adaptador

1. Añadir valor al enum `AIAdapters`
2. Implementar adaptador con `AIServiceInterface`
3. Registrar en `adaptersMap` de `AIServiceFactory`

---

## 💾 Cache Services

**Ubicación**: `src/infrastructure/services/cache/`

Servicios de caché en memoria para reducir consultas a la base de datos.

### Cachés Disponibles

#### 1. embeddingCache

```typescript
import { embeddingCache } from "@/infrastructure/services/cache/InMemoryEmbeddingCache";

const embeddings = embeddingCache.getEmbeddings();
embeddingCache.setEmbeddings(newEmbeddings);
```

#### 2. aiEntityCache

```typescript
import { aiEntityCache } from "@/infrastructure/services/cache/InMemoryAiEntityCache";

const aiEntities = aiEntityCache.getAiEntities();
aiEntityCache.setAiEntities(newAiEntities);
```

#### 3. operationInProgressRegistry

Previene ejecución duplicada de operaciones costosas.

```typescript
import { operationInProgressRegistry } from "@/infrastructure/services/cache/OperationInProgressRegistry";
import { UpdateHandlerType } from "@domain/interfaces/cache/OperationInProgressRegistryInterface";

// Verificar si está en progreso
const isInProgress = operationInProgressRegistry.has(
    pageId,
    statusId,
    UpdateHandlerType.PAGE,
);

// Registrar operación
const promise = someAsyncOperation();
operationInProgressRegistry.set(pageId, statusId, UpdateHandlerType.PAGE, promise);

// Limpiar al completar
operationInProgressRegistry.delete(pageId, statusId, UpdateHandlerType.PAGE);
```

### Directrices

- **Siempre** usa los singleton proporcionados
- Los datos persisten hasta reinicio del servidor
- Nuevas instancias no comparten estado

---

## 🛡️ Resilience Services

**Ubicación**: `src/domain/services/resilience/`

Servicios que implementan patrones de resiliencia para manejar fallos elegantemente.

### Circuit Breaker

**Archivo**: `src/domain/services/resilience/CircuitBreaker.ts`

Previene fallos en cascada deteniendo peticiones a servicios que están fallando.

#### Estados

- **CLOSED**: Operación normal, peticiones pasan
- **OPEN**: Umbral de fallos alcanzado, peticiones fallan rápido
- **HALF_OPEN**: Probando recuperación

#### Uso

```typescript
const breaker = new CircuitBreaker({
    threshold: 3,
    resetTimeout: 30000,
    name: "ExternalAPI",
});

if (breaker.isOpen()) {
    return; // Omitir operación
}

try {
    await expensiveOperation();
    breaker.recordSuccess();
} catch (error) {
    breaker.recordFailure();
    throw error;
}
```

### Database Query Limiter

**Archivo**: `src/domain/services/resilience/DatabaseQueryLimiter.ts`

Controla la carga de la base de datos limitando consultas concurrentes.

```typescript
const limiter = new DatabaseQueryLimiter({
    maxConcurrentQueries: 15,
    name: "PrimaryDB",
});

const result = await limiter.execute(async () => {
    return await userRepository.find();
});

// Verificar presión
if (limiter.isUnderPressure()) {
    console.warn("La cola está creciendo");
}
```

### Mejores Prácticas

- Establecer umbrales apropiados basados en confiabilidad
- Usar nombres significativos para monitoreo
- Monitorear cambios de estado
- Combinar con reintentos para fallos transitorios

---

## 🔐 Sistema de Permisos

**Ubicación**: `src/domain/permissions/`

Sistema de permisos flexible y tipado para controlar acceso a funcionalidades.

### Componentes

#### 1. permissions.constants.ts

Define todas las constantes de permisos (autogenerado).

```typescript
import { GriddoPermissions } from "@domain/permissions/permissions.constants";

const permiso = GriddoPermissions.GENERAL_ACCESS_TO_SITES; // 'general.accessToSites'
```

#### 2. permissions.types.ts

Tipos TypeScript para validación.

```typescript
import type { PermissionKey } from "@domain/permissions/permissions.types";

function checkPermission(permissionKey: PermissionKey) {
    // TypeScript valida que la clave sea válida
}

checkPermission("general.accessToSites"); // ✅
checkPermission("permiso.inexistente"); // ❌ Error de TypeScript
```

### Categorías de Permisos

1. **general**: Acceso y gestión de sitios
2. **usersRoles**: Gestión de usuarios y roles
3. **content**: Páginas y contenido
4. **seoAnalytics**: SEO y analíticas
5. **mediaGallery**: Gestión de medios
6. **navigation**: Gestión de navegación
7. **categories**: Gestión de categorías
8. **forms**: Gestión de formularios
9. **global**: Permisos a nivel global

### Directrices

> **Importante**: Los permisos se generan automáticamente mediante script. No modifiques manualmente `permissions.constants.ts`.

Para añadir nuevos permisos:
1. Actualiza la configuración del script `generate-permissions.ts`
2. Ejecuta el script para regenerar
3. Usa la convención de nomenclatura existente

---

## 🔄 Lifecycle Services

**Ubicación**: `src/infrastructure/services/lifecycle/`

Servicios para gestionar el ciclo de vida de la aplicación.

### Inicialización de Servicios

El sistema requiere que varios servicios se inicialicen al arrancar:

```typescript
// En el arranque de la aplicación
await CommandBusService.initialize();   // Comandos síncronos
await EventBusService.initialize();     // Eventos asíncronos
await QueueWorkerService.initialize();  // Colas y jobs
await QueueWorkerService.startProcessing(); // Iniciar procesamiento
```

### Orden de Inicialización

1. Base de datos (`DBTypeORM`)
2. CommandBus
3. EventBus
4. QueueWorkerService
5. Cachés (lazy loading)

### Graceful Shutdown

Asegura que los recursos se liberen correctamente al cerrar:

```typescript
process.on("SIGTERM", async () => {
    await QueueWorkerService.stopProcessing();
    await db.close();
    process.exit(0);
});
```

---

## 📚 Referencias Rápidas

### Comandos de Desarrollo

```bash
# Tests
yarn test                  # Todos los tests
yarn test:watch           # Modo watch
yarn test:coverage        # Con cobertura
yarn test:integration     # Solo integración

# Base de datos
yarn migration:generate    # Generar migración
yarn migration:run        # Ejecutar migraciones
yarn migration:revert     # Revertir última
yarn migration:show       # Ver estado

# Build
yarn build                # Construir proyecto
yarn dev                  # Modo desarrollo
```

### Patrones de Inyección de Dependencias

```typescript
// En Casos de Uso
constructor(
    private readonly repository: SomeRepositoryInterface,
    private readonly eventBus: EventBusInterface,
    private readonly commandBus: CommandBusInterface,
) {}

// En Controladores
const repository = new SomeRepository(db, new UuidGenerator());
const useCase = new SomeUseCase(repository, eventBus, commandBus);
```

### Errores Comunes

| Error | Causa | Solución |
|-------|-------|----------|
| `Cannot find module` | Path alias incorrecto | Usa `@domain`, `@application`, etc. |
| `Repository method not implemented` | Interfaz sin implementar | Implementa método o elimina de interfaz |
| `Database not initialized` | Setup de tests falló | Revisa `vitest.setup.ts` |
| `Zod validation error` | Datos inválidos | Revisa esquema Zod en `src/infrastructure/dto/` |

---

## ✅ Checklist de Calidad para PRs
- [ ] ¿Se han creado/actualizado los tipos en `api-types` si afectan al cliente?
- [ ] ¿El endpoint tiene su esquema Zod de entrada y salida?
- [ ] ¿Está el endpoint registrado y documentado en Swagger?
- [ ] ¿La lógica de negocio está en el Caso de Uso y no en el Controlador?
- [ ] ¿Se inyectan las dependencias por interfaz?
- [ ] ¿Existen tests unitarios para el Caso de Uso?
- [ ] ¿Existen tests de integración para la ruta?
- [ ] ¿Se usan los servicios apropiados (CommandBus/EventBus/Queue)?
- [ ] ¿Se usa `HttpResponseHandler` y `DatabaseConnector` correctamente?

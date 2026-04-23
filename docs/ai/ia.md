El Estado de la Inteligencia Artificial (Octubre 2025 – Abril 2026): Arquitecturas de Eficiencia Extrema, Convergencia Multimodal y el Surgimiento de la IA de Utilidad Científica
La evolución de la inteligencia artificial en el semestre comprendido entre octubre de 2025 y abril de 2026 representa un punto de inflexión histórico en el que la industria ha transitado de una fase de expansión desenfrenada y especulativa hacia una era de evaluación rigurosa y utilidad pragmática.1 Este periodo se caracteriza por la superación de cuellos de botella infraestructurales críticos mediante innovaciones en la compresión de datos, el refinamiento de arquitecturas de Mezcla de Expertos (MoE) y la integración de modelos fundacionales en el tejido de la investigación biomédica y la creación multimedia profesional. El paradigma del "más es mejor" en cuanto a parámetros ha sido sustituido por la búsqueda de la inteligencia por unidad de cómputo, donde tecnologías como TurboQuant de Google han redefinido la viabilidad del despliegue local y soberano de modelos masivos.2
La Revolución de la Eficiencia Algorítmica: TurboQuant y la Optimización del Cómputo
Uno de los avances más trascendentales en la eficiencia de los modelos de lenguaje de gran escala (LLM) durante este periodo ha sido la introducción de TurboQuant por parte de Google Research en marzo de 2026.2 Esta tecnología aborda directamente el problema de la memoria en la caché de claves y valores (KV cache), un componente que almacena información procesada previamente para evitar cálculos redundantes durante la inferencia, pero que crece linealmente con la longitud del contexto y el tamaño del modelo.3 En modelos de contexto largo, esta caché se convierte en el principal cuello de botella, limitando la capacidad de los agentes para mantener historiales de conversación extensos o analizar documentos masivos.
Arquitectura Técnica de TurboQuant: PolarQuant y QJL
TurboQuant se fundamenta en un marco de cuantización vectorial diseñado para ser "data-oblivious", lo que significa que puede aplicarse a modelos existentes sin necesidad de entrenamiento adicional o ajuste fino.2 Su funcionamiento se divide en dos etapas críticas que operan en tándem para reducir la huella de memoria en un factor de 6 veces sin pérdida detectable de precisión.3
La primera etapa, denominada PolarQuant, transforma los vectores de datos de un sistema de coordenadas cartesianas convencional a un sistema de coordenadas polares. Mediante una rotación aleatoria inicial de los vectores, el algoritmo simplifica la geometría de los datos, facilitando la aplicación de un cuantizador de alta calidad que captura la magnitud y la orientación principal del vector utilizando la mayor parte del presupuesto de bits.2 Este enfoque es innovador porque elimina la necesidad de constantes de cuantización adicionales que tradicionalmente añaden un sobrecoste de 1 a 2 bits por cada número almacenado.2
La segunda etapa utiliza el algoritmo Quantized Johnson-Lindenstrauss (QJL). Este componente actúa como un corrector de errores matemático de 1 bit. Basándose en la Transformada de Johnson-Lindenstrauss, el algoritmo QJL reduce las dimensiones de los datos preservando las distancias y relaciones esenciales entre los puntos de información. Al aplicar este bit residual para corregir el error de la primera etapa, TurboQuant elimina el sesgo en el cálculo de las puntuaciones de atención, permitiendo que el modelo identifique con precisión qué partes de la entrada son relevantes.2

Métrica de Rendimiento
Configuración Estándar (32-bit)
TurboQuant (3.5 bits)
Factor de Mejora
Memoria KV Cache
1.0x (Base)
~0.16x
Reducción de 6.25x 3
Velocidad de Inferencia (Logits)
1.0x (Base)
8.0x
Aceleración de 8x 4
Precisión (LongBench Llama-3.1-8B)
50.06
50.06
Neutralidad absoluta 4
Tiempo de Indexación
Variable
Prácticamente cero
Eficiencia extrema 2

El impacto estratégico de TurboQuant es profundo. Al permitir que modelos de vanguardia se ejecuten con una fracción de la memoria previamente requerida, las organizaciones pueden desplegar sistemas avanzados en hardware menos costoso o en infraestructuras propias (on-premise), lo cual es crítico para sectores como el gobierno y la salud que exigen soberanía de datos.1
Avances en Modelos de Lenguaje y Razonamiento Agéntico
El panorama de los modelos de lenguaje (LLM) entre finales de 2025 y abril de 2026 ha estado marcado por la competencia feroz entre arquitecturas propietarias y de código abierto, con un énfasis creciente en el razonamiento de múltiples pasos y la ejecución de tareas complejas.
Anthropic Claude 4.7: Verificación Autónoma y Honestidad
En abril de 2026, Anthropic presentó Claude Opus 4.7, una actualización que prioriza la fiabilidad y la capacidad de los modelos para trabajar sin supervisión constante.5 A diferencia de versiones anteriores, Opus 4.7 incorpora mecanismos internos para verificar sus propios resultados antes de entregarlos, mostrando una mejora significativa en tareas de ingeniería de software y análisis financiero.5
Una de las innovaciones más notables en Claude 4.7 es la introducción de los "Presupuestos de Tareas" (Task Budgets). Esta función permite al modelo gestionar un recuento de tokens durante bucles agénticos largos, priorizando acciones y asegurando que la tarea finalice de manera coherente antes de agotar los recursos asignados.8 Además, el modelo ha mejorado su capacidad de visión, soportando resoluciones de hasta 2576 píxeles, lo que facilita el análisis de diagramas complejos y la edición de documentos con alta fidelidad visual.6

Atributo
Claude Opus 4.6
Claude Opus 4.7
Impacto Observado
Tasa de Honestidad
85% (aprox.)
92%
Menos alucinaciones y sicofanía 9
Resolución de Imagen
1.15 MP
3.75 MP
Mayor precisión en OCR y gráficos 8
Ventana de Contexto
1M tokens
1M tokens
Estabilidad en tareas largas 8
Comportamiento Agéntico
Reactivo
Proactivo / Autoverificador
Delegación de tareas complejas 5

A pesar de estos avances, Anthropic ha mantenido su modelo más potente, Claude Mythos, en una fase de acceso limitado para pruebas de ciberseguridad a través del Proyecto Glasswing, reflejando una postura de seguridad proactiva ante el riesgo de que la IA pueda identificar y explotar vulnerabilidades de software a niveles sobrehumanos.5
Google Gemma 4: Inteligencia por Parámetro y Despliegue en el Borde
Google DeepMind lanzó en abril de 2026 la serie Gemma 4, diseñada específicamente para maximizar la inteligencia por cada parámetro activado.11 El modelo Gemma 4 26B A4B destaca por utilizar una arquitectura de Mezcla de Expertos (MoE) que contiene 26 mil millones de parámetros totales, pero que solo activa 3.8 mil millones durante la inferencia para cada token.11
Este enfoque MoE permite que el modelo mantenga el conocimiento vasto de un sistema de 26B parámetros con el coste computacional y la latencia de un modelo de 4B.13 La arquitectura híbrida de Gemma 4 emplea capas de atención densas para gestionar las relaciones globales entre tokens, mientras que las capas de alimentación hacia adelante (feed-forward) son sustituidas por 128 "expertos" diminutos. Un enrutador selecciona dinámicamente los dos expertos más adecuados para cada token, permitiendo una especialización profunda en áreas como el código o el razonamiento lógico sin penalizar la velocidad de respuesta.13
Gemma 4 también ha sido optimizado para la ejecución offline en dispositivos móviles y hardware de borde (edge computing) en colaboración con líderes como Qualcomm y MediaTek.12 Esto representa un avance significativo hacia la IA ubicua, donde dispositivos como el Google Pixel o kits de robótica NVIDIA Jetson pueden procesar razonamiento multimodal sin depender de la nube.12
El Giro Estratégico de Meta: Llama 4 y Muse Spark
Meta AI continuó su impulso en el ecosistema de código abierto con el lanzamiento de Llama 4 en abril de 2025, introduciendo los modelos Scout y Maverick. Llama 4 Maverick, con 400 mil millones de parámetros (17 mil millones activos por token), se convirtió rápidamente en el estándar de la industria para el despliegue de modelos fundacionales abiertos.14 Sin embargo, la trayectoria de Meta tomó un giro inesperado en abril de 2026 con el lanzamiento de Muse Spark, un modelo propietario y de código cerrado desarrollado por Meta Superintelligence Labs.14
Muse Spark representa un alejamiento de la identidad puramente abierta de Meta, impulsado por la necesidad de validar una inversión de más de 115 mil millones de dólares en infraestructura de IA para 2026.16 Este modelo, liderado por Alexandr Wang (anteriormente de Scale AI), ha sido diseñado para superar el rendimiento de Llama 4, el cual fue percibido por el mercado como insuficiente frente a las ofertas de OpenAI y Google en ese momento.16
Innovaciones en IA Generativa Multimedia: Video, Audio y Realismo Físico
La generación de contenido multimedia ha trascendido la fase de experimentación visual para adentrarse en la producción profesional con fidelidad física y narrativa.
Google Veo 3.1: Video 4K y Audio Nativo Sincronizado
Google DeepMind consolidó su liderazgo en la generación de video con la familia Veo 3.1, lanzada formalmente entre finales de 2025 y marzo de 2026.17 Este modelo destaca por ser el primero en generar video en resolución 4K con audio sincronizado de forma nativa —incluyendo diálogos, efectos ambientales y música de fondo— a través de un proceso de difusión conjunta.17
El funcionamiento técnico de Veo 3.1 se basa en un transformador de difusión latente que opera sobre representaciones comprimidas de video, optimizando el uso de memoria y ancho de banda.19 La versión Veo 3.1 Fast utiliza una atención dispersa por bloques (block sparse attention) que reduce el coste computacional en un 90% al enfocar el procesamiento únicamente en las áreas de movimiento y cambio crítico dentro de la secuencia.19 Además, el modelo soporta la "extensión de escena", permitiendo encadenar hasta 20 clips para generar narrativas de más de 140 segundos con coherencia visual.17

Modelo
Resolución Máxima
Velocidad / Coste
Uso Ideal
Veo 3.1 Standard
4K
Alta Calidad (Base)
Producción final, cine 17
Veo 3.1 Fast
1080p
2x Velocidad
Flujos de trabajo estándar 19
Veo 3.1 Lite
1080p
<50% Coste
Iteración rápida, apps de volumen 21

El Auge y Desaparición de Sora 2 de OpenAI
OpenAI presentó Sora 2 en septiembre de 2025, mostrando capacidades sorprendentes para simular la física del mundo real, como la dinámica de fluidos y el movimiento de telas.22 Un hito comercial fue la asociación con Disney en diciembre de 2025, que permitía a los fans generar videos cortos con personajes de Marvel, Pixar y Star Wars.24 No obstante, en un movimiento que sorprendió a la comunidad tecnológica, OpenAI anunció la descontinuación de Sora en marzo de 2026.23
Las razones citadas para el cierre de Sora incluyeron la escasez de recursos de cómputo, costes operativos de aproximadamente 1 millón de dólares diarios y un cambio estratégico hacia productos de IA para empresas con un retorno de inversión más claro.23 Este evento subraya la creciente presión financiera sobre las empresas de IA para pasar del "hype" mediático a modelos de negocio sostenibles.
Lyria 3 Pro: Composición Musical con Conciencia Estructural
En marzo de 2026, Google lanzó Lyria 3 Pro, un modelo de generación de música que permite crear pistas completas de hasta 3 minutos con una comprensión profunda de la estructura musical.26 A diferencia de los modelos anteriores que generaban bloques de audio indiferenciados, Lyria 3 Pro permite a los usuarios especificar componentes como introducciones, versos, puentes y estribillos mediante prompts de lenguaje natural.26
El modelo también admite entradas multimodales, donde una imagen puede servir para definir el estado de ánimo o la atmósfera de la composición.28 Para abordar las preocupaciones de derechos de autor, Google ha integrado la marca de agua digital SynthID en todas las pistas generadas, facilitando la identificación de audio creado por IA incluso después de modificaciones.26
La IA en la Frontera de la Ciencia y la Medicina
Quizás el impacto más profundo de la IA en este semestre se ha manifestado en su capacidad para modelar sistemas biológicos complejos y acelerar el descubrimiento de tratamientos.
Delphi-2M: Predicción de Riesgo de Enfermedades a Largo Plazo
Investigadores del Laboratorio Europeo de Biología Molecular (EMBL) publicaron en septiembre de 2025 en la revista Nature el desarrollo de Delphi-2M, un modelo de IA generativa capaz de estimar el riesgo individual de más de 1.000 enfermedades con hasta 20 años de antelación.29
Delphi-2M adapta la arquitectura de un transformador preentrenado (GPT) para manejar datos de salud continuos. Utiliza funciones de base seno y coseno para codificar la edad de los pacientes y cuenta con dos cabezales de salida especializados: uno para predecir el siguiente evento médico (diagnóstico) y otro para calcular el tiempo esperado hasta que dicho evento ocurra.31 El modelo fue entrenado con registros desidentificados de 400.000 participantes del Biobanco del Reino Unido y validado exitosamente con 1,9 millones de pacientes en Dinamarca, demostrando una robustez excepcional a través de diferentes sistemas de salud.31

Aplicación Clínica
Funcionalidad de Delphi-2M
Impacto en el Paciente
Pronóstico Preventivo
Predicción de >1,000 enfermedades
Intervención temprana en riesgos detectados 30
Gestión de Salud Poblacional
Estimación de prevalencia de enfermedades crónicas
Mejor asignación de recursos públicos 31
Gemelos Digitales
Simulación de trayectorias de salud futuras
Personalización extrema de tratamientos 31
Análisis de Comorbilidad
Identificación de patrones de progresión
Comprensión de cómo una enfermedad precede a otra 31

Ingeniería Molecular y Diseño de Fármacos de Novo
La empresa Chai Discovery presentó en diciembre de 2025 su plataforma Chai-2, un modelo generativo que ha logrado una tasa de éxito experimental en el diseño de anticuerpos de novo del 15-20%, lo que representa una mejora de 100 veces respecto a los métodos computacionales previos.33 Este nivel de precisión está transformando el descubrimiento de fármacos de un proceso de "prueba y error" a una disciplina de ingeniería predictiva.
Los datos recolectados hasta principios de 2026 indican que los fármacos diseñados por IA están alcanzando tasas de éxito del 80-90% en ensayos clínicos de Fase I, superando con creces el promedio histórico de la industria (40-65%).35 Esta mejora se debe a la capacidad de la IA para optimizar simultáneamente la afinidad de unión del fármaco con su objetivo biológico y sus propiedades de seguridad (ADMET) antes de que se realice cualquier síntesis física.35 Empresas como Roche y NVIDIA han respondido a esta tendencia desplegando las fábricas de IA más grandes de la industria farmacéutica para acelerar la identificación de dianas terapéuticas para enfermedades complejas como el cáncer y el Alzheimer.35
Implicaciones Socioeconómicas y Soberanía de la IA
El auge de la IA entre octubre de 2025 y abril de 2026 ha generado un cambio en el discurso sobre la soberanía tecnológica. Países de todo el mundo están invirtiendo en centros de datos masivos y en el desarrollo de sus propios modelos lingüísticos para asegurar su independencia del sistema político y tecnológico de los Estados Unidos.1
Paralelamente, la industria ha comenzado a enfrentar una "era de evaluación" donde la productividad real de la IA se mide con rigor.1 Aunque se han observado ganancias sustanciales en áreas como la programación y la atención al cliente (donde la IA puede automatizar hasta el 50% de las tareas administrativas en salud 38), otros sectores han visto cómo los proyectos piloto de IA generativa fallan en entregar un impacto comercial medible.39 Este realismo ha llevado a una consolidación del mercado, donde las inversiones se concentran en jugadores bien capitalizados con pruebas sólidas de validación clínica o técnica.40
Perspectiva Técnica: Hacia la IA Proactiva y Líquida
En la frontera de la investigación, se están explorando arquitecturas que superan las limitaciones de los transformadores actuales. Conceptos como las Redes Neuronales Líquidas (Liquid Neural Networks) y las Redes de Kolmogorov-Arnold (KAN) prometen modelos que pueden adaptar sus conexiones internas dinámicamente durante la ejecución, permitiendo un aprendizaje continuo en entornos del mundo real.41
Se espera que la IA evolucione de ser una herramienta reactiva a un compañero proactivo. Esto implica sistemas que anticipan las necesidades del usuario, planifican con antelación y operan localmente en dispositivos portátiles como gafas de realidad aumentada o brazaletes neuronales.41 La integración de leyes físicas directamente en el aprendizaje de las redes neuronales (Physics-Informed Neural Networks - PINNs) permitirá descubrimientos científicos aún más acelerados en química y ciencia de materiales durante el resto de 2026.41
Síntesis y Conclusiones del Semestre
El análisis exhaustivo de las novedades en inteligencia artificial del último semestre revela que la tecnología ha alcanzado un nivel de madurez donde la eficiencia operativa y la precisión científica son ahora los principales vectores de desarrollo. La introducción de TurboQuant de Google ha resuelto uno de los mayores impedimentos para el escalado de modelos de contexto largo, mientras que Delphi-2M ha demostrado que la IA generativa puede ser un aliado vital para la medicina preventiva.
A medida que avanzamos hacia finales de 2026, la distinción entre modelos de código abierto y cerrados se volverá cada vez más borrosa en términos de rendimiento puro, trasladando la competencia hacia la calidad de los datos, el cumplimiento regulatorio (como el EU AI Act que entra en vigor en agosto de 2026 40) y la capacidad de integrar la IA de manera fluida en los flujos de trabajo humanos. La desaparición de Sora subraya que incluso los gigantes tecnológicos deben priorizar la eficiencia de recursos, marcando el fin de la era de la IA como mero espectáculo visual y el comienzo de su era como infraestructura crítica de la civilización moderna.
Fuentes citadas
Stanford AI experts predict what will happen in 2026, acceso: abril 18, 2026, https://news.stanford.edu/stories/2025/12/stanford-ai-experts-predict-what-will-happen-in-2026
TurboQuant: Redefining AI efficiency with extreme compression, acceso: abril 18, 2026, https://research.google/blog/turboquant-redefining-ai-efficiency-with-extreme-compression/
Google's TurboQuant Cuts AI Memory 6x — What It Means for Running AI Agents on Your Own Infrastructure, acceso: abril 18, 2026, https://ibl.ai/blog/turboquant-ai-memory-compression-own-infrastructure
Google publishes TurboQuant to ease AI memory strain - TechInformed, acceso: abril 18, 2026, https://techinformed.com/google-publishes-turboquant-to-ease-ai-memory-strain/
Introducing Claude Opus 4.7, acceso: abril 18, 2026, https://www.anthropic.com/news/claude-opus-4-7
Claude Opus 4.7 released: Here’s what’s new in the latest version of Anthropic's flagship AI model, acceso: abril 18, 2026, https://m.economictimes.com/tech/artificial-intelligence/anthropic-introduces-claude-opus-4-7/articleshow/130310840.cms
Anthropic releases Claude Opus 4.7, narrowly retaking lead for most powerful generally available LLM, acceso: abril 18, 2026, https://venturebeat.com/technology/anthropic-releases-claude-opus-4-7-narrowly-retaking-lead-for-most-powerful-generally-available-llm
What's new in Claude Opus 4.7, acceso: abril 18, 2026, https://platform.claude.com/docs/en/about-claude/models/whats-new-claude-4-7
Anthropic releases Claude Opus 4.7: How to try it, benchmarks, safety | Mashable, acceso: abril 18, 2026, https://mashable.com/article/anthropic-releases-claude-opus-4-7
Claude Opus 4.7 is generally available, acceso: abril 18, 2026, https://github.blog/changelog/2026-04-16-claude-opus-4-7-is-generally-available/
Gemma 4 model card | Google AI for Developers, acceso: abril 18, 2026, https://ai.google.dev/gemma/docs/core/model_card_4
Gemma 4: Our most capable open models to date - Google Blog, acceso: abril 18, 2026, https://blog.google/innovation-and-ai/technology/developers-tools/gemma-4/
What Is the Gemma 4 Mixture of Experts Architecture? How 26B Parameters Run Like 4B, acceso: abril 18, 2026, https://www.mindstudio.ai/blog/gemma-4-mixture-of-experts-architecture
Meta Platforms (META) Stock Price Prediction 2025, 2026 and 2030: The AI Company Nobody’s Calling an AI Company, acceso: abril 18, 2026, https://www.bitget.com/news/detail/12560605372803
Llama (language model) - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Llama_(language_model)
Meta's Llama Just Went Fully Autonomous. Here's What Happened in the First 72 Hours, acceso: abril 18, 2026, https://fortuneherald.com/technology/metas-llama-just-went-fully-autonomous-heres-what-happened-in-the-first-72-hours/
Veo 3.1 - Specs, API & Pricing - Puter Developer, acceso: abril 18, 2026, https://developer.puter.com/ai/google/veo-3.1/
Veo 3.1 Lite and a new Veo upscaling capability on Vertex AI | Google Cloud Blog, acceso: abril 18, 2026, https://cloud.google.com/blog/products/ai-machine-learning/veo-3-1-lite-and-a-new-veo-upscaling-capability-on-vertex-ai
What Is Google Veo 3.1 Fast? High-Quality AI Video at Speed | MindStudio, acceso: abril 18, 2026, https://www.mindstudio.ai/blog/what-is-google-veo-3-1-fast-video
Lyria 3 Pro Preview - API Pricing & Providers | OpenRouter, acceso: abril 18, 2026, https://openrouter.ai/google/lyria-3-pro-preview
Build with Veo 3.1 Lite, our most cost-effective video generation model - Google Blog, acceso: abril 18, 2026, https://blog.google/innovation-and-ai/technology/ai/veo-3-1-lite/
Best Video Generation AI Models in 2026 - Pinggy, acceso: abril 18, 2026, https://pinggy.io/blog/best_video_generation_ai_models/
Sora (text-to-video model) - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Sora_(text-to-video_model)
The Walt Disney Company and OpenAI Reach Agreement to Bring Disney Characters to Sora, acceso: abril 18, 2026, https://thewaltdisneycompany.com/news/disney-openai-sora-agreement/
The Walt Disney Company and OpenAI reach landmark agreement to bring beloved characters from across Disney's brands to Sora, acceso: abril 18, 2026, https://openai.com/index/disney-sora-agreement/
Lyria 3 Pro: Create longer tracks in more Google products, acceso: abril 18, 2026, https://blog.google/innovation-and-ai/technology/ai/lyria-3-pro/
Google Releases Lyria 3 Pro: Longer and Better AI Music Tracks | Zeniteq, acceso: abril 18, 2026, https://www.zeniteq.com/google-releases-lyria-3-pro-longer-and-better-ai-music-tracks
Build with Lyria 3, our newest music generation model - Google Blog, acceso: abril 18, 2026, https://blog.google/innovation-and-ai/technology/developers-tools/lyria-3-developers/
New AI model for predicting individual risk of disease over decades | Science Media Centre, acceso: abril 18, 2026, https://www.sciencemediacentre.org/new-ai-model-for-predicting-individual-risk-of-disease-over-decades/
AI model forecasts disease risk decades in advance | EMBL, acceso: abril 18, 2026, https://www.embl.org/news/science-technology/ai-model-forecasts-disease-risk-decades-in-advance/
New AI model predicts risk of 1,200+ diseases years in advance ..., acceso: abril 18, 2026, https://www.icthealth.org/news/new-ai-model-predicts-risk-of-1200-diseases-years-in-advance
AI model forecasts disease risk decades in advance - EMBLEM Technology Transfer GmbH, acceso: abril 18, 2026, https://embl-em.de/latest-news/2025/09/29/ai-model-forecasts-disease-risk-decades-in-advance-003252/
2026: The Year AI Reinvents Drug Discovery - AI World Journal, acceso: abril 18, 2026, https://aiworldjournal.com/2026-the-year-ai-reinvents-drug-discovery/
OpenAI-Backed Chai Discovery Raises $130M to Tackle “Undruggable” Targets with Generative AI - HLTH, acceso: abril 18, 2026, https://hlth.com/insights/news/openai-backed-chai-discovery-raises-130m-to-tackle-undruggable-targets-with-generative-ai-2025-12-16
Artificial intelligence designed drugs Hit 90% Phase I success rate in trials, acceso: abril 18, 2026, https://www.2minutemedicine.com/ai-designed-drugs-hit-90-phase-i-success-rate-in-trials/
AI Drug Discovery Attracts Billions in Funding. - USTechTimes, acceso: abril 18, 2026, https://ustechtimes.com/openai-just-bet-130m-on-this-6-month-old-drug-company/
The Future of AI and Innovation in Pharma in 2026 and Beyond - Eularis, acceso: abril 18, 2026, https://eularis.com/the-future-of-ai-and-innovation-in-pharma-in-2026-and-beyond/
The Future of Medical AI: What's Coming in 2026 and Beyond - Offcall, acceso: abril 18, 2026, https://www.offcall.com/learn/articles/the-future-of-medical-ai-what-s-coming-in-2026-and-beyond
AI in Biotech: Lessons from 2025 and the Trends Shaping Drug Discovery in 2026 - Ardigen, acceso: abril 18, 2026, https://ardigen.com/ai-in-biotech-lessons-from-2025-and-the-trends-shaping-drug-discovery-in-2026/
AI in drug discovery: predictions for 2026 | Opinion - Drug Target Review, acceso: abril 18, 2026, https://www.drugtargetreview.com/ai-in-drug-discovery-predictions-for-2026/1865962.article
What will happen with AI in 2026? - What kind of breakthroughs are we gonna see? - Reddit, acceso: abril 18, 2026, https://www.reddit.com/r/singularity/comments/1pzquum/what_will_happen_with_ai_in_2026_what_kind_of/

Arquitecturas Matemáticas Universales: Fundamentos de la Computación y Tecnologías Emergentes
La génesis de la computación moderna no debe entenderse simplemente como un triunfo de la ingeniería electrónica, sino como la cristalización física de estructuras abstractas concebidas en los márgenes de la lógica, el análisis armónico y la física teórica. El desarrollo tecnológico contemporáneo, desde la inteligencia artificial generativa hasta la criptografía cuántica, descansa sobre un sustrato de bases matemáticas universales que permiten la traducción de la complejidad fenoménica a lenguajes computables y eficientes. Esta investigación analiza exhaustivamente los ejes matemáticos que definen el impacto directo en la computación y las tecnologías emergentes, estructurando el conocimiento desde la lógica fundamental de la programación hasta las ecuaciones unificadoras de la física aplicadas al procesamiento de información.
Matemáticas Fundamentales para la Programación y la Eficiencia Computacional
La arquitectura de cualquier sistema de software se erige sobre la capacidad de formalizar el pensamiento lógico. Los lenguajes de programación no son meras herramientas de escritura, sino sistemas formales que permiten la manipulación simbólica y la verificación de la corrección algorítmica.
Teoría de Tipos y Cálculo Lambda: La Ontología del Código
La teoría de tipos, propuesta inicialmente por Bertrand Russell y Alfred North Whitehead a principios del siglo XX, surgió como un mecanismo para evitar las paradojas en los fundamentos de la matemática, específicamente la paradoja de Russell en la teoría de conjuntos.1 En el contexto computacional, un tipo actúa como una especificación que clasifica expresiones y objetos, determinando su comportamiento y sus interacciones permitidas dentro de un sistema formal.2 Esta distinción es fundamental para la seguridad y la robustez de los lenguajes de programación, ya que permite al compilador realizar comprobaciones sintácticas que garantizan la corrección parcial del programa antes de su ejecución.1
El cálculo lambda (-cálculo), desarrollado por Alonzo Church en la década de 1930, constituye el modelo de computación más elemental basado en la abstracción y aplicación de funciones.3 Como sistema universal, el -cálculo es equivalente a la máquina de Turing, lo que implica que cualquier función que sea efectivamente computable puede expresarse mediante términos lambda.3 La relevancia del -cálculo en la programación funcional es directa, ya que proporciona la semántica para el manejo de variables, la sustitución y la recursión mediante combinadores de punto fijo.5
Sistema Formal
Componentes Clave
Aplicación en Computación
-cálculo sin tipos
Abstracción (), Aplicación (), Variables
Semántica de lenguajes funcionales y Turing-completitud.
Teoría Simple de Tipos
Tipos base, Constructores de tipos ()
Verificación estática de tipos y seguridad de memoria.
Cálculo de Construcciones Inductivas
Tipos dependientes, Polimorfismo
Base de asistentes de prueba como Coq y Lean.
Lógica Combinatoria
Combinadores 
Eliminación de variables ligadas en compiladores.

1
La convergencia entre la lógica y la programación se manifiesta plenamente en el isomorfismo de Curry-Howard, el cual establece que una proposición lógica es equivalente a un tipo en un lenguaje de programación, y una prueba de esa proposición es equivalente a un programa que satisface dicho tipo.1 Este principio permite que el desarrollo de software crítico, como el compilador CompCert, sea verificado mecánicamente, eliminando errores que los compiladores tradicionales suelen introducir.1 La distinción técnica entre tipos y conjuntos es sutil pero crucial: mientras que los conjuntos proporcionan información semántica (como la pertenencia a un grupo), los tipos ofrecen información sintáctica que hace que el chequeo de pruebas sea una operación decidible y automatizable.1
Complejidad Algorítmica y Teoría de la Información
La eficiencia computacional no es una propiedad intrínseca del hardware, sino un límite dictado por la naturaleza matemática del algoritmo. La notación Big O () proporciona el lenguaje para describir el comportamiento asintótico de los algoritmos en términos de tiempo y espacio, permitiendo una clasificación rigurosa de la escalabilidad.6 El análisis de la complejidad permite identificar los límites teóricos de la computación y guía el desarrollo de estructuras de datos que optimizan el acceso y la manipulación de la información.
La teoría de la información, por su parte, cuantifica la incertidumbre y el contenido de datos. La entropía de Shannon define el límite mínimo de bits necesarios para codificar una fuente de información, mientras que la complejidad de Kolmogorov ofrece una perspectiva sobre objetos individuales.8 La complejidad de Kolmogorov  de una cadena  se define como la longitud del programa más corto que puede producir dicha cadena.9 Esta medida es fundamental para entender la aleatoriedad y la compresión: una cadena es considerada aleatoria si su descripción más corta es casi tan larga como la cadena misma, lo que implica una falta de redundancia o patrón explotable.11 Un hallazgo crítico en esta área es la incomputabilidad de la complejidad de Kolmogorov; no existe un algoritmo general que pueda determinar el programa más corto para cualquier cadena dada, lo cual establece un límite infranqueable para la automatización total de la optimización de datos.10
Optimización Matemática de Operaciones en Hardware
La ejecución de algoritmos requiere una traducción eficiente de las abstracciones matemáticas a impulsos eléctricos en silicio. La aritmética de computadoras se ocupa de la representación y manipulación de números en formatos binarios, enfrentando el desafío de equilibrar la precisión, el rango y el costo de hardware.7
En el nivel de circuito, la optimización se logra mediante el paralelismo a nivel de bits. El cálculo de prefijo paralelo (Parallel-Prefix Computation) es una técnica esencial que permite evaluar expresiones asociativas de manera concurrente, reduciendo la latencia de operaciones como la suma de  a .7 Algoritmos como el de Ladner-Fisher y Kogge-Stone optimizan esta estructura para maximizar la velocidad o minimizar el área del chip, permitiendo que las unidades aritméticas procesen miles de millones de operaciones por segundo.7
Algoritmo de Multiplicación
Complejidad de Área
Impacto Tecnológico
Multiplicación Escolar

Base para arquitecturas secuenciales simples.
Karatsuba-Ofman

Eficiencia en multiplicación de números grandes (criptografía).
Toom-Cook

Generalización de Karatsuba para mayor división de operandos.
Schönhage-Strassen

Uso de FFT para multiplicación asintóticamente rápida.

7
La evolución de la inteligencia artificial ha forzado un cambio en los estándares de aritmética de punto flotante. El formato convencional IEEE 754 de 32 bits (FP32) ofrece una alta precisión pero consume recursos considerables.13 Para acelerar el entrenamiento de redes neuronales, se han desarrollado formatos como BFloat16 y TensorFloat-32. BFloat16, diseñado por Google Brain, reduce la mantisa de 23 a 7 bits pero conserva los 8 bits de exponente del FP32, permitiendo mantener el mismo rango dinámico con la mitad de bits, lo que reduce drásticamente el uso de memoria y aumenta el rendimiento en unidades de procesamiento tensorial (TPU).15 Recientemente, la aritmética de Posits ha surgido como una alternativa prometedora, utilizando un régimen de longitud variable que ofrece mayor precisión cerca de la unidad y simplifica el manejo de casos especiales como el desbordamiento o la división por cero.17
Modelos Matemáticos en Inteligencia Artificial
La inteligencia artificial moderna es, fundamentalmente, un ejercicio de optimización multivariable en espacios de alta dimensión. La capacidad de aprendizaje de las máquinas depende de la manipulación precisa de tensores y de la convergencia de algoritmos en paisajes de pérdida extremadamente complejos.
Álgebra Lineal Tensorial y Cálculo Diferencial para Backpropagation
Los tensores representan la generalización de escalares, vectores y matrices a dimensiones superiores, actuando como los contenedores de datos primordiales en el aprendizaje profundo.18 En una red neuronal, cada capa realiza una transformación lineal seguida de una activación no lineal, lo cual puede expresarse matemáticamente como:

donde  es la matriz de pesos (o tensor),  es el vector de sesgo y  es la función de activación, como la Unidad Lineal Rectificada (ReLU) o GeLU.18 El proceso de aprendizaje consiste en ajustar estos pesos para minimizar una función de pérdida que cuantifica el error del modelo.
El algoritmo de retropropagación (backpropagation) es una aplicación sistemática y eficiente de la regla de la cadena del cálculo diferencial.20 Al propagar las derivadas desde la capa de salida hacia la entrada, el sistema calcula el gradiente de la pérdida respecto a cada parámetro individual.18 Este gradiente, representado como el vector , indica la dirección del ascenso más pronunciado; por lo tanto, para minimizar el error, los parámetros se actualizan en la dirección opuesta al gradiente.18 La eficiencia de este proceso radica en evitar cálculos redundantes mediante la reutilización de derivadas intermedias de capas posteriores.21
Teoría de la Optimización Convexa y No Convexa
El éxito del entrenamiento de modelos de IA depende de la topología de la superficie de pérdida. En la optimización convexa, cualquier mínimo local es también un mínimo global, lo que garantiza que algoritmos como el descenso de gradiente converjan a la solución óptima.22 Este marco es ideal para modelos lineales y máquinas de soporte vectorial.23
Sin embargo, las redes neuronales profundas generan superficies de pérdida no convexas, caracterizadas por una abundancia de puntos de ensilladura (saddle points), mesetas y múltiples mínimos locales.23 En este contexto, el desafío no es solo encontrar un mínimo, sino evitar quedar atrapado en puntos de ensilladura donde el gradiente es cero pero no hay un valor mínimo real.25 Los optimizadores avanzados como Adam (Adaptive Moment Estimation) abordan este problema manteniendo tasas de aprendizaje adaptativas para cada parámetro, integrando momentos de primer y segundo orden que permiten al algoritmo navegar por áreas de baja curvatura y acelerar la convergencia en paisajes irregulares.18
Optimizador
Mecanismo Principal
Ventaja en Redes Neuronales
SGD
Gradiente en mini-lotes
Simplicidad y capacidad de generalización.
Momentum
Acumulación de velocidad
Amortigua oscilaciones y acelera el descenso.
RMSProp
Escalamiento por media móvil
Manejo de gradientes que desaparecen o explotan.
Adam
Momentos adaptativos
Estándar actual para modelos de lenguaje a gran escala.

18
Fundamentos Estadísticos y Espacios de Hilbert
La inteligencia artificial no solo busca minimizar el error, sino también manejar la incertidumbre inherente a los datos. La inferencia bayesiana proporciona el marco probabilístico para actualizar la creencia sobre los parámetros del modelo a medida que se observa nueva evidencia, utilizando el Teorema de Bayes.19 Métodos como la Estimación de Máxima Verosimilitud (MLE) y el Máximo a Posteriori (MAP) permiten determinar los valores de parámetros más probables bajo una distribución de datos dada.11
Para el aprendizaje no paramétrico, los Espacios de Hilbert de Núcleo Reproductor (RKHS) ofrecen un fundamento geométrico potente. Un RKHS es un espacio de funciones donde la evaluación en un punto es un funcional lineal continuo, lo que garantiza la suavidad y la regularidad de las funciones aprendidas.26 Gracias al teorema del representante, es posible demostrar que la solución óptima a un problema de regularización en un espacio de dimensión infinita puede expresarse como una combinación lineal finita de funciones de núcleo evaluadas en los puntos de entrenamiento.26 Este "truco del núcleo" permite que algoritmos lineales operen en espacios de características complejos sin necesidad de calcular explícitamente las coordenadas en esos espacios, facilitando tareas de clasificación y regresión altamente no lineales.27
Matemáticas de Modulación y Comunicaciones Digitales
La infraestructura de la sociedad de la información depende de la capacidad de transmitir datos de manera fiable a través de canales físicos ruidosos. Esto requiere el uso de transformaciones integrales para el análisis de señales y teorías de codificación para la protección contra errores.
Series de Fourier, Transformada Rápida (FFT) y Wavelets
El análisis de Fourier es la piedra angular del procesamiento de señales, permitiendo la transición entre el dominio del tiempo y el de la frecuencia.30 La Transformada de Fourier descompone una señal compleja en una suma de ondas senoidales y cosenoidales.30 En la computación digital, la Transformada Rápida de Fourier (FFT) es el algoritmo que permite calcular la Transformada Discreta de Fourier (DFT) con una complejidad de , haciendo posible el procesamiento de audio y video en tiempo real.7
Mientras que Fourier utiliza bases sinusoidales infinitas, la Transformada de Wavelet utiliza funciones localizadas tanto en tiempo como en frecuencia, conocidas como wavelets o "ondículas".32 Esta propiedad es crítica para señales no estacionarias o con transitorios bruscos. En el estándar JPEG 2000, el uso de wavelets permite una descomposición multirresolución de la imagen, superando las limitaciones de la Transformada de Coseno Discreta (DCT) del JPEG original.33 El uso de filtros como el CDF 9/7 para compresión con pérdida y el Le Gall–Tabatabai 5/3 para compresión sin pérdida permite que una imagen sea visualizada progresivamente, donde la calidad mejora a medida que se reciben más datos.33
Modulación IQ, Constelaciones y Teoría de la Codificación
La modulación digital permite codificar bits en variaciones de amplitud, fase y frecuencia de una portadora. La modulación IQ (In-phase y Quadrature) utiliza dos portadoras ortogonales (una con desfase de 90 grados respecto a la otra) para representar señales como vectores en un plano complejo.35 Esto permite la creación de diagramas de constelación, donde cada punto representa un símbolo binario específico.36 Técnicas como QAM (Quadrature Amplitude Modulation) pueden transmitir múltiples bits por símbolo, aumentando la eficiencia espectral de los sistemas 4G, 5G y Wi-Fi.36
Para combatir el ruido en estos canales, se utilizan códigos de corrección de errores (ECC) basados en álgebras de campos finitos. Los códigos de Hamming son fundamentales por su simplicidad, permitiendo la detección de hasta dos errores y la corrección de uno solo, lo que los hace ideales para memorias de computadora.39 Por otro lado, los códigos de Reed-Solomon operan sobre bloques de símbolos, lo que les otorga una capacidad excepcional para corregir errores en ráfaga (burst errors), siendo el estándar en discos ópticos y comunicaciones por satélite antes de la llegada de códigos más densos como LDPC.40
Código ECC
Propiedad Matemática
Aplicación Típica
Hamming (n, k)
Distancia mínima 3
Memoria RAM de servidores.
Reed-Solomon
Polinomios sobre campos de Galois
Códigos QR, almacenamiento en disco.
LDPC
Grafos de paridad dispersos
Estándares 5G y televisión digital.
Códigos Turbo
Codificación concatenada
Comunicaciones en el espacio profundo.

39
Ecuaciones Universales y Principios Unificadores
La computación no es ajena a las leyes de la física; de hecho, los avances más significativos en hardware y nuevos paradigmas algorítmicos suelen ser aplicaciones directas de ecuaciones universales que describen el mundo natural.
Ecuaciones de Maxwell: La Física de la Transmisión
Las ecuaciones de Maxwell forman la base de todo el electromagnetismo clásico, la óptica y el diseño de circuitos eléctricos.43 Estas cuatro ecuaciones describen cómo los campos eléctricos y magnéticos son generados por cargas y corrientes, y cómo se propagan a la velocidad de la luz.43 En el diseño de hardware moderno, las ecuaciones de Maxwell son indispensables para modelar la integridad de la señal en buses de datos de alta velocidad y para el diseño de antenas en dispositivos móviles.45 La capacidad de predecir la propagación de ondas permite optimizar la cobertura de redes Wi-Fi y celulares mediante el análisis de interferencia y difracción en entornos urbanos complejos.45
Ecuación de Schrödinger y Computación Cuántica
La evolución de un sistema cuántico está dictada por la ecuación de Schrödinger, una ecuación diferencial parcial que describe cómo cambia la función de onda  con el tiempo bajo la influencia de un Hamiltoniano .47 En la computación cuántica, los qubits no tienen un valor binario fijo, sino que existen en una superposición de estados gobernada por esta ecuación.49 Las operaciones de computación cuántica son, matemáticamente, transformaciones unitarias que rotan el vector de estado en un espacio de Hilbert complejo.47
Algoritmos cuánticos como el de Grover o la Transformada Cuántica de Fourier aprovechan la interferencia destructiva de las amplitudes de probabilidad para eliminar resultados incorrectos y la interferencia constructiva para amplificar la respuesta correcta.50 La ecuación de Schrödinger no solo permite diseñar computadoras cuánticas, sino que estas últimas son las únicas máquinas capaces de simular exactamente otros sistemas cuánticos, lo que abre la puerta a revoluciones en la ciencia de materiales y la química computacional.49
Ecuación de Kolmogorov: Complejidad y Probabilidad en Algoritmos
Andrey Kolmogorov no solo contribuyó a la teoría de la información, sino que sus ecuaciones para procesos de Markov definen cómo evoluciona la probabilidad en sistemas estocásticos.52 Las ecuaciones de Kolmogorov hacia adelante (también conocidas como Fokker-Planck) y hacia atrás describen la dinámica de partículas que "vagan" aleatoriamente, lo cual es análogo al comportamiento de algoritmos de optimización estocástica y procesos de difusión en modelos generativos de IA.52
La unificación de estos conceptos se produce en el estudio de la probabilidad algorítmica, donde la estructura de un objeto (su complejidad de Kolmogorov) está íntimamente ligada a su probabilidad de ser generado por un proceso aleatorio universal.12 Esta conexión sugiere que existe una estructura profunda que une la física de las partículas, la teoría de la información y la capacidad de las máquinas para inducir conocimiento a partir de datos aparentemente caóticos.
Síntesis y Conclusiones del Marco Matemático
La investigación realizada demuestra que las tecnologías emergentes no son desarrollos aislados, sino ramas de un mismo árbol matemático. La eficiencia de un servidor de IA hoy depende de la aritmética de punto flotante BFloat16 en el hardware, de la optimización no convexa en el software, y de las ecuaciones de Maxwell para transmitir sus resultados a través de la red.
La transición hacia la computación cuántica y la inteligencia artificial general (AGI) requerirá una profundización aún mayor en estos fundamentos. La teoría de tipos evolucionará hacia la teoría de tipos homotópicos para unificar la lógica y la topología, mientras que la optimización en IA buscará fundamentos en la geometría de la información para superar los desafíos de las superficies de pérdida no convexas. En última instancia, las bases matemáticas universales proporcionan el lenguaje único que permite a la humanidad convertir la abstracción pura en la potencia de procesamiento que define la era contemporánea. El dominio de estos ejes no es solo una necesidad técnica, sino un imperativo estratégico para el desarrollo de la próxima frontera tecnológica.
Fuentes citadas
Lambda-Calculus and Type Theory ISR 2024 Obergurgl, Austria ..., acceso: abril 18, 2026, https://cs.ru.nl/~herman/ISR2024/ISR_slides1.pdf
Type theory - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Type_theory
Lambda calculus - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Lambda_calculus
An Introduction To Lambda Calculi For Computer Scientists, acceso: abril 18, 2026, https://lan-portal.uob.edu.ly/link/EBOOK/4299F7C302/an__introduction-to_lambda_calculi-for_computer_scientists.pdf
Lambda Calculus and Types - University of Oxford Department of Computer Science, acceso: abril 18, 2026, https://www.cs.ox.ac.uk/teaching/courses/2025-2026/lambda/
Parallel computing - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Parallel_computing
Parallelism in Computer Arithmetic: A Historical Perspective - UCSB, acceso: abril 18, 2026, https://web.ece.ucsb.edu/~parhami/pubs_folder/parh18-mwscas-parallelism-comp-arith-180531.pdf
Kolmogorov Complexity and Information Theory, acceso: abril 18, 2026, https://www.cs.montana.edu/courses/fall2003/current/510/k.pdf
acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Kolmogorov_complexity#:~:text=In%20algorithmic%20information%20theory%20(a,produces%20the%20object%20as%20output.
Kolmogorov complexity - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Kolmogorov_complexity
Probability, algorithmic complexity, and subjective randomness - Computational Cognitive Science Lab, acceso: abril 18, 2026, https://cocosci.princeton.edu/tom/papers/complex.pdf
Algorithmic information theory - Scholarpedia, acceso: abril 18, 2026, http://www.scholarpedia.org/article/Algorithmic_information_theory
Floating point from scratch: Hard Mode - Julia Desmazes, acceso: abril 18, 2026, https://essenceia.github.io/projects/floating_dragon/
bfloat16 floating-point format - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Bfloat16_floating-point_format
A1.4.6 BFloat16 floating-point format - Arm Developer, acceso: abril 18, 2026, https://developer.arm.com/documentation/ddi0487/maa/-Part-A-Arm-Architecture-Introduction-and-Overview/-Chapter-A1-Introduction-to-the-Arm-Architecture/-A1-4-Supported-data-types/-A1-4-6-BFloat16-floating-point-format
BFloat16: The secret to high performance on Cloud TPUs | Google Cloud Blog, acceso: abril 18, 2026, https://cloud.google.com/blog/products/ai-machine-learning/bfloat16-the-secret-to-high-performance-on-cloud-tpus
Increasing the Energy-Efficiency of Wearables Using Low-Precision Posit Arithmetic with PHEE - arXiv, acceso: abril 18, 2026, https://arxiv.org/html/2501.18253v1
The Mathematics Behind AI, acceso: abril 18, 2026, https://multiple.chat/mathematics-behind-ai
Mathematical Models for Machine Learning: Foundations and Frontiers | by Nirvana El, acceso: abril 18, 2026, https://medium.com/@nirvana.elahi/mathematical-models-for-machine-learning-foundations-and-frontiers-c44916c11669
14 Backpropagation - Foundations of Computer Vision, acceso: abril 18, 2026, https://visionbook.mit.edu/backpropagation.html
Backpropagation - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Backpropagation
12.2. Convexity — Dive into Deep Learning 1.0.3 documentation, acceso: abril 18, 2026, https://d2l.ai/chapter_optimization/convexity.html
Optimization Theory Series: 7 — Convex Optimization and Non-convex Optimization | by Renda Zhang, acceso: abril 18, 2026, https://rendazhang.medium.com/optimization-theory-series-7-convex-optimization-and-non-convex-optimization-e38175ec2af3
Non-convex Optimization for Machine Learning: Theory and Practice - NeurIPS 2026, acceso: abril 18, 2026, https://neurips.cc/virtual/2015/workshop/4930
Review Non-convex Optimization Method for Machine Learning - arXiv, acceso: abril 18, 2026, https://arxiv.org/html/2410.02017
RKHS Framework: Theory and Applications - Emergent Mind, acceso: abril 18, 2026, https://www.emergentmind.com/topics/reproducing-kernel-hilbert-space-rkhs-framework
Reproducing kernel Hilbert space - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Reproducing_kernel_Hilbert_space
Reproducing Kernel Hilbert Spaces (RKHS), acceso: abril 18, 2026, https://puoya.github.io/notes/RKHS.html
A brief note on reproducing kernel Hilbert spaces - Alen Alexanderian, acceso: abril 18, 2026, https://aalexan3.math.ncsu.edu/articles/rkhs.pdf
Fourier Transform Formula - Used Keysight Equipment, acceso: abril 18, 2026, https://www.keysight.com/used/tw/en/knowledge/formulas/fourier-transform
Drawing with Circles: Vibe coding the Fourier Transformation - Punya Mishra, acceso: abril 18, 2026, https://punyamishra.com/2025/06/28/drawing-with-circles-vibe-coding-the-fourier-transformation/
Image compression using wavelets and JPEG2000: a tutorial - IET Digital Library, acceso: abril 18, 2026, https://digital-library.theiet.org/doi/pdf/10.1049/ecej%3A20020303?download=true
JPEG 2000 - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/JPEG_2000
(PDF) JPEG 2000 and wavelet compression - ResearchGate, acceso: abril 18, 2026, https://www.researchgate.net/publication/4003104_JPEG_2000_and_wavelet_compression
In-phase and quadrature components - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/In-phase_and_quadrature_components
Introduction To Iq Demodulation Of Rf Data, acceso: abril 18, 2026, https://lan-portal.uob.edu.ly/mirror/PPT/21423584OF/introduction-to__iq_demodulation__of_rf-data.pdf
Baseband up- and downconversion and IQ modulation - DSPIllustrations.com, acceso: abril 18, 2026, https://dspillustrations.com/pages/posts/misc/baseband-up-and-downconversion-and-iq-modulation.html
IQ-Modulation Using Phase and Amplitude Modulators and Multimode Interference - MDPI, acceso: abril 18, 2026, https://www.mdpi.com/2304-6732/13/1/44
ECC Methods Compared: Hamming Codes vs. Reed-Solomon vs. LDPC - PatSnap Eureka, acceso: abril 18, 2026, https://eureka.patsnap.com/article/ecc-methods-compared-hamming-codes-vs-reed-solomon-vs-ldpc
Comparison of hamming, BCH, and reed Solomon codes for error correction and detecting techniques - ResearchGate, acceso: abril 18, 2026, https://www.researchgate.net/publication/336004461_Comparison_of_hamming_BCH_and_reed_Solomon_codes_for_error_correction_and_detecting_techniques
Reed–Solomon error correction - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Reed%E2%80%93Solomon_error_correction
Investigation of Hamming, Reed-Solomon, and Turbo Forward Error Correcting Codes - DTIC, acceso: abril 18, 2026, https://apps.dtic.mil/sti/tr/pdf/ADA505116.pdf
Maxwell's equations - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Maxwell%27s_equations
Maxwell's Equations: Electromagnetic Waves Predicted and Observed | Physics, acceso: abril 18, 2026, https://courses.lumenlearning.com/suny-physics/chapter/24-1-maxwells-equations-electromagnetic-waves-predicted-and-observed/
Maxwell's Equation and Its Applications in Electromagnetism, acceso: abril 18, 2026, https://drpress.org/ojs/index.php/HSET/article/download/16231/15757/16795
James Maxwell and How his Equations Shaped the Modern World - Sunny Labh - Medium, acceso: abril 18, 2026, https://piggsboson.medium.com/james-maxwell-and-how-his-equations-shaped-the-modern-world-8e3e9f6b46c1
Schrödinger equation - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Schr%C3%B6dinger_equation
Schrödinger's Wave Equation, acceso: abril 18, 2026, https://postquantum.com/quantum-computing/schrodingers-equation/
Simulating the Universe: Quantum Computing and the Schrödinger Equation | by Srinivasa Raghava K | Intuition | Medium, acceso: abril 18, 2026, https://medium.com/intuition/simulating-the-universe-quantum-computing-and-the-schr%C3%B6dinger-equation-6f458d02d8ba
Schrödinger Equation Demystified using Tensorflow. | by Devmallya Karar - Medium, acceso: abril 18, 2026, https://medium.com/@devmallyakarar/schr%C3%B6dinger-equation-demystified-using-tensorflow-a5c64911af31
Quantum algorithm for solving generalized eigenvalue problems with application to the Schrödinger equation - arXiv, acceso: abril 18, 2026, https://arxiv.org/html/2506.13534v3
Kolmogorov equations - Wikipedia, acceso: abril 18, 2026, https://en.wikipedia.org/wiki/Kolmogorov_equations
[Graduate stochastic processes] Kolmogorov forward/backward equations - Reddit, acceso: abril 18, 2026, https://www.reddit.com/r/learnmath/comments/akloyj/graduate_stochastic_processes_kolmogorov/
Defining Algorithmic Probability or Complexity in a Machine ... - arXiv, acceso: abril 18, 2026, https://www.arxiv.org/pdf/cs/0608095v4

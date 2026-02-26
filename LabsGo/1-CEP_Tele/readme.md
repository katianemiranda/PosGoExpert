Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

Requisitos - Serviço A (responsável pelo input):

O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
Caso seja válido, será encaminhado para o Serviço B via HTTP
Caso não seja válido, deve retornar:
Código HTTP: 422
Mensagem: invalid zipcode
Requisitos - Serviço B (responsável pela orquestração):

O sistema deve receber um CEP válido de 8 digitos
O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
O sistema deve responder adequadamente nos seguintes cenários:
Em caso de sucesso:
Código HTTP: 200
Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
Em caso de falha, caso o CEP não seja válido (com formato correto):
Código HTTP: 422
Mensagem: invalid zipcode
​​​Em caso de falha, caso o CEP não seja encontrado:
Código HTTP: 404
Mensagem: can not find zipcode
Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

Implementar tracing distribuído entre Serviço A - Serviço B
Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura
Dicas:

Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/
Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
Sendo F = Fahrenheit
Sendo C = Celsius
Sendo K = Kelvin
Para dúvidas da implementação do OTEL, você pode clicar aqui
Para implementação de spans, você pode clicar aqui
Você precisará utilizar um serviço de collector do OTEL
Para mais informações sobre Zipkin, você pode clicar aqui
Entrega:

O código-fonte completo da implementação.
Documentação explicando como rodar o projeto em ambiente dev.
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.

*** Observações sobre o projeto ***

Este projeto implementa dois serviços em Go (Service A e Service B) para consultar a temperatura com base em um CEP, aplicando Tracing Distribuído com OpenTelemetry (OTEL) e Zipkin.

O projeto é organizado em uma estrutura de monorepo simplificada com subdiretórios para cada serviço.

o Zipkin é o repositório central onde todos os dados de rastreamento (traces) dos serviços A e B são armazenados.

O Service A utiliza otelhttp.NewHandler e otelhttp.DefaultClient juntamente com propagation.TraceContext{} para garantir que o trace context seja injetado nos headers da requisição HTTP (W3C Trace Context) e, assim, o Service B consiga continuar o mesmo trace.

No Service B, foi criado spans específicos (getCityFromCEP, getTemperature) para isolar e medir o tempo de resposta de cada API externa, como solicitado.

Tratamento de Erros: O código lida com os três cenários de erro solicitados: formato inválido (422), CEP não encontrado (404) e erros internos/APIs (500).

## Estrutura do Projeto

- **service-a**: Responsável por receber o input (`{ "cep": "..." }`), validar e encaminhar a requisição para o Service B. Implementa instrumentação OTEL para requisições de entrada e saída.
- **service-b**: Responsável pela orquestração, buscando a cidade (ViaCEP) e a temperatura (WeatherAPI), realizando as conversões (C, F, K) e respondendo. Implementa Spans para medir as latências das chamadas externas.
- **docker-compose.yml**: Define os três serviços (Service A, Service B, Zipkin).

*** PARA EXECUÇÃO VEJA O ARQUIVO DOCUMENTO.MD 
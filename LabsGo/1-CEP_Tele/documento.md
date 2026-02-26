### EXECUÇÃO

docker-compose up --build

Os serviços estarão rodando nas seguintes portas:

Serviço	Porta	Função
Service A	localhost:8080	Endpoint Principal
Service B	localhost:8081	Serviço Interno
Zipkin UI	localhost:9411	Visualizador de Traces

Testar a Aplicação
Use o curl ou um cliente HTTP (como Postman/Insomnia) para enviar uma requisição para o Service A na porta 8080.

Exemplo de Sucesso (CEP Válido)

curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -d '{"cep": "01001000"}'

Response (HTTP 200):

{
    "city": "São Paulo",
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.5
}

Exemplo de Falha (CEP Inválido - Formato)

curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -d '{"cep": "123"}'

Response (HTTP 422):

{
    "message": "invalid zipcode"
}

Exemplo de Falha (CEP Não Encontrado)

curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -d '{"cep": "99999999"}'

Response (HTTP 404):

{
    "message": "can not find zipcode"
}

Visualizar o Tracing no Zipkin

Acesse a interface do Zipkin em: http://localhost:9411

Você verá traces completos que mostram o fluxo:

Service A (Incoming Request)

Service A (processCEPHandler - Span pai)

Service A -> Requisição HTTP para Service B (propagação de contexto)

Service B (getWeatherHandler - Span pai)

Service B (getCityFromCEP - Span para ViaCEP)

Service B (getTemperature - Span para WeatherAPI)

Os spans 5 e 6 irão medir o tempo de resposta exato de cada API externa.


Exemplo de como usar o Zipkin

acesse: http://localhost:9411.
Aba "Find Traces" (Encontrar Traces), você pode filtrar por:
Nome do Serviço: (Ex: service-a ou service-b).
Nome do Span: (Ex: getCityFromCEP ou getWeatherHandler).
Tags/Anotações: (Ex: buscar todos os traces com a tag http.status_code=404).
Visualizar Dados: Ao clicar em um trace, você vê todos os spans, seus tempos de latência e as tags (como CEP, cidade, e códigos de erro).
Exportação: O Zipkin permite salvar a visualização do trace, mas se você precisar de exportação programática, siga o método abaixo.
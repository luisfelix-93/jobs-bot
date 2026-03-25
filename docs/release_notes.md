# Release Notes

## Correções
* **API de Vagas (Findwork):** O processamento de vagas da integração com a API Findwork foi corrigido. O ID retornado estava sofrendo falhas de conversão de tipo. A tipagem do identificador de vagas foi ajustada internamente de número inteiro (`int`) para texto (`string`), garantindo assim a correta ingestão no sistema através da modelagem de domínio (`GUID`).

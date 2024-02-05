<div align="center">
<img height="150px" src="https://raw.githubusercontent.com/wisdom-oss/brand/main/svg/standalone_color.svg">
<h1>Smartmeter Data Management</h1>
<h3>service-smartmeter-rest</h3>
<p>💧 data management for smart meters</p>
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-smartmeter-rest?style=for-the-badge" alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="OpenAPI Schema Version"/>
</a>
<a href="https://github.com/wisdom-oss/service-smartmeter-rest/pkgs/container/service-smartmeter-rest">
<img alt="Static Badge" src="https://img.shields.io/badge/ghcr.io-wisdom--oss%2Fservice--smartmeter--rest-2496ED?style=for-the-badge&logo=docker&logoColor=white&labelColor=555555">
</a>
</div>

> [!NOTE]  
> This service does not interact with smart meters and their configuration as
> this should only be done by authorized personnel.

> [!IMPORTANT]
> This microservice requires a PostgreSQL database with the 
> [timescaleDB](https://github.com/timescale/timescaledb)
> extension installed.

This microservice allows the management of smart meter data that has been 
written into the database.
Furthermore, the service provides an endpoint allowing external services to
write collected smart meter data into the database, either as a batch operation
or writing single entries.
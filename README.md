# LPA Agent

## Quick Start

Write to `config.json`

```json
{
  "program": "./lpac",
  "env": {
    "APDU_INTERFACE": "./libapduinterface_pcsc.so",
    "HTTP_INTERFACE": "./libhttpinterface_curl.so"
  }
}
```

Run the program

```bash
./lpa-agent -c config.json
```

## References

- [SGP.21 v2.5](https://www.gsma.com/esim/resources/sgp-21-v2-5/)
- [SGP.22 v2.5](https://www.gsma.com/esim/resources/sgp-22-v2-5/)

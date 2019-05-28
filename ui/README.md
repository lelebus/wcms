## Prerequisites

### Yarn
Install **yarn** following the instructions [here](https://yarnpkg.com/en/docs/install)

## Development

- From the terminal, navigate to this folder
- Install **parcel-bundler** `yarn global add parcel-bundler`
- Run `parcel src/index.html`

## Deployment

Build the files into _static/_ folder `parcel build src/index.html && cp -r dist/ ../static/`

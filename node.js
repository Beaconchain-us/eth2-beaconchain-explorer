// config.js
const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');

const env = process.env.NODE_ENV || 'local'; // local, staging, production

let configFile;
switch(env) {
  case 'production':
    configFile = 'production.yml';
    break;
  case 'staging':
    configFile = 'staging.yml';
    break;
  default:
    configFile = 'local.yml';
}

const configPath = path.join(__dirname, configFile);
const config = yaml.load(fs.readFileSync(configPath, 'utf8'));

console.log(`✅ Loaded config from ${configFile} (${env} mode)`);
module.exports = config;
import { mkdtempSync, writeFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { execFileSync, spawn } from 'node:child_process';

const root = mkdtempSync(join(tmpdir(), 'shiyu-e2e-'));
const config = join(root, 'config.yaml');
writeFileSync(config, `data_dir: ${JSON.stringify(root)}
server:
  port: "18808"
  mode: release
database:
  driver: sqlite
  database: ${JSON.stringify(join(root, 'app.db'))}
jwt:
  secret: ""
  expire_time: 3600
  issuer: shiyu-e2e
log:
  level: warn
`);
writeFileSync('.p0-e2e.json', JSON.stringify({ root }));
execFileSync(process.env.GO_BIN || 'go', ['build', '-o', join(root, 'server'), './cmd/server'], { cwd: resolve('../../backend/shiyu-admin-backend'), stdio: 'inherit' });
const child = spawn(join(root, 'server'), [], { env: { ...process.env, CONFIG_FILE: config }, stdio: 'inherit' });
for (const signal of ['SIGINT', 'SIGTERM']) process.on(signal, () => child.kill(signal));
child.on('exit', (code) => { rmSync(root, { recursive: true, force: true }); rmSync('.p0-e2e.json', { force: true }); process.exit(code || 0); });

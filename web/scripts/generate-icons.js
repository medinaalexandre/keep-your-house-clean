import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const publicDir = path.join(__dirname, '../public');
const svgPath = path.join(publicDir, 'favicon.svg');

async function generateIcons() {
  try {
    const sharp = await import('sharp');
    
    const svgBuffer = fs.readFileSync(svgPath);
    
    const sizes = [192, 512];
    
    for (const size of sizes) {
      await sharp.default(svgBuffer)
        .resize(size, size)
        .png()
        .toFile(path.join(publicDir, `icon-${size}.png`));
      
      console.log(`✓ Gerado icon-${size}.png`);
    }
    
    console.log('\n✅ Todos os ícones foram gerados com sucesso!');
  } catch (error) {
    if (error.code === 'MODULE_NOT_FOUND' && error.message.includes('sharp')) {
      console.error('❌ Erro: A biblioteca "sharp" não está instalada.');
      console.error('\nPara instalar, execute:');
      console.error('  cd web && npm install --save-dev sharp');
      process.exit(1);
    } else {
      console.error('❌ Erro ao gerar ícones:', error.message);
      process.exit(1);
    }
  }
}

generateIcons();


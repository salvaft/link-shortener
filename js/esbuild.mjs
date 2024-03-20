import * as esbuild from "esbuild";
import { globSync } from "glob";

const tsFiles = globSync("src/**/*.ts");

await esbuild.build({
  entryPoints: tsFiles,
  format: "esm",
  platform: "browser",
  minify: false,
  bundle: true,
  outdir: "../public/js/",
});

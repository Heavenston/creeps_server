import * as esbuild from "esbuild";
import { tailwindPlugin } from 'esbuild-plugin-tailwindcss';
import { glob } from "glob";
import { $ } from "zx";

await $`mkdir -p dist`;

await Promise.all([
  $`cp -r static/* dist`,
  esbuild.build({
    entryPoints: await glob("src/entrypoints/*.ts"),
    outdir: "dist",
    bundle: true,

    format: "esm",

    splitting: true,
    minify: true,
    treeShaking: true,
    sourcemap: "inline",

    tsconfig: "tsconfig.json",

    define: Object.fromEntries(
      Object.entries(process.env).map(([k, v]) => [`process.env.${k}`, JSON.stringify(v)])
    ),
  }),
  esbuild.build({
    entryPoints: await glob("src/**/*.css"),
    outdir: "dist",

    minify: true,
    treeShaking: true,
    sourcemap: "inline",

    plugins: [
      tailwindPlugin({ }),
    ],
  }),
]);

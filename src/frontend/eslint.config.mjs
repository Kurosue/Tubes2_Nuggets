import path from "path";
import { fileURLToPath } from "url";

import { includeIgnoreFile } from "@eslint/compat";
import { FlatCompat } from "@eslint/eslintrc";
import javascript from "@eslint/js";
import typescript from "typescript-eslint";
import stylistic from "@stylistic/eslint-plugin";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const gitignoreFile = path.join(__dirname, ".gitignore");
const tsConfigFile = path.join(__dirname, "tsconfig.json");
const compat = new FlatCompat({ baseDirectory: __dirname });

/** @type {import("eslint").Linter.Config[]} */
export default [
	includeIgnoreFile(gitignoreFile),
	javascript.configs.recommended,
	...typescript.configs.recommended.map(c => ({ ...c, files: [...(c.files ?? []), "**/*.ts", "**/*.tsx"] })),
	stylistic.configs.recommended,
	...compat.extends("next"),
	{
		files: [
			"**/*.ts",
			"**/*.tsx"
		],
		rules: {
			"@typescript-eslint/no-explicit-any": "off",
			"@typescript-eslint/no-unused-vars": "off",
			"@typescript-eslint/no-useless-constructor": "off",
			"@typescript-eslint/no-unnecessary-type-assertion": "error",
			"@typescript-eslint/no-unnecessary-boolean-literal-compare": "error",
			"@typescript-eslint/no-non-null-asserted-optional-chain": "error",
			"@typescript-eslint/no-non-null-asserted-nullish-coalescing": "error",
			"@typescript-eslint/non-nullable-type-assertion-style": "error",
			"@typescript-eslint/strict-boolean-expressions": ["error", {
				allowString: false,
				allowNumber: false,
				allowNullableObject: false
			}]
		}
	},
	{
		languageOptions: {
			parser: typescript.parser,
			parserOptions: {
				project: tsConfigFile,
				tsconfigRootDir: __dirname
			}
		},
		rules: {
			"no-empty-pattern": "off",
			"curly": ["warn", "multi-or-nest"],
			"eqeqeq": "off",
			"import/no-anonymous-default-export": "off",
			"no-bitwise": "off",
			"no-return-assign": "off",
			"no-unused-vars": "off",
			"no-useless-constructor": "off",
			"no-empty": ["error", { "allowEmptyCatch": true }],
			"no-constant-condition": ["warn", { "checkLoops": false }],
			"object-shorthand": "off",
			"require-await": "off",
			"@next/next/no-html-link-for-pages": ["error", "app/"],
			"@stylistic/no-trailing-spaces": ["warn"],
			"@stylistic/arrow-parens": ["warn", "as-needed"],
			"@stylistic/semi": ["warn", "always"],
			"@stylistic/comma-dangle": ["warn", "never"],
			"@stylistic/indent": ["warn", "tab", { SwitchCase: 1 }],
			"@stylistic/jsx-closing-bracket-location": ["warn", "line-aligned"],
			"@stylistic/jsx-indent": ["warn", "tab"],
			"@stylistic/jsx-indent-props": ["warn", "tab"],
			"@stylistic/jsx-quotes": ["warn", "prefer-double"],
			"@stylistic/jsx-self-closing-comp": "off",
			"@stylistic/keyword-spacing": ["warn", {
				after: true,
				overrides: {
					if: { after: false },
					for: { after: false },
					while: { after: false },
					catch: { after: false },
					switch: { after: false },
					await: { after: false }
				}
			}],
			"@stylistic/no-tabs": "off",
			"@stylistic/object-curly-spacing": ["warn", "always"],
			"@stylistic/quotes": ["warn", "double", { avoidEscape: true, allowTemplateLiterals: true }],
			"@stylistic/member-delimiter-style": ["warn", {
				multilineDetection: "brackets",
				multiline: {
					delimiter: "semi",
					requireLast: true
				},
				singleline: {
					delimiter: "comma",
					requireLast: false
				}
			}],
			"@stylistic/indent-binary-ops": ["warn", "tab"],
			"@stylistic/multiline-ternary": "off",
			"@stylistic/operator-linebreak": ["warn", "after"],
			"@stylistic/brace-style": ["warn", "1tbs", { allowSingleLine: true }],
			"@stylistic/max-statements-per-line": ["warn", { max: 5 }],
			"@stylistic/lines-between-class-members": "off",
			"@stylistic/space-before-function-paren": "off",
			"@stylistic/jsx-one-expression-per-line": "off",
			"@stylistic/quote-props": "off",
			"@stylistic/no-multi-spaces": ["warn", {
				ignoreEOLComments: true
			}]
		}
	}
];

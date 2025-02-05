/*eslint-env node*/
// import franken from 'franken-ui/shadcn-ui/preset-quick';

/** @type {import('tailwindcss').Config} */
export default {
	// presets: [franken()],
	content: ["./views/*.{html,js,templ}", "./views/**/*.{html,js,templ}"],
	safelist: [
		{
			pattern: /^uk-/
		},
		'ProseMirror',
		'ProseMirror-focused',
		'tiptap'
	],
	theme: {
		extend: {}
	},
	plugins: []
};

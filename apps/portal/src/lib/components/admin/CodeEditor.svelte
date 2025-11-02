<!--
  A reusable Svelte component for the Monaco Editor (the engine for VS Code).
  - Handles SSR by only loading the editor on the client-side.
  - Configured for Python syntax highlighting.
  - Uses a custom 'Monokai Light' theme.
  - Enables bracket pair colorization ('rainbow brackets').
  - Supports two-way binding for the editor's content.
-->
<script lang="ts">
	import { onMount } from 'svelte';
	import type * as Monaco from 'monaco-editor';

	let { value = $bindable('') } = $props<{ value: string }>();

	let editorEl: HTMLDivElement;
	let editor = $state<Monaco.editor.IStandaloneCodeEditor | undefined>(undefined);

	let isUpdatingFromParent = false;

	onMount(() => {
		async function initializeMonaco() {
			const monaco = await import('monaco-editor');

			monaco.editor.defineTheme('monokai-light', {
				base: 'vs',
				inherit: true,
				rules: [
					{ token: 'keyword', foreground: '8959a8' },
					{ token: 'string', foreground: '718c00' },
					{ token: 'number', foreground: 'f5871f' },
					{ token: 'comment', foreground: '8e908c', fontStyle: 'italic' },
					{ token: 'predefined', foreground: 'c82829' },
					{ token: 'operator', foreground: '3e999f' },
					{ token: 'identifier', foreground: '4271ae' },
					{ token: 'delimiter', foreground: '3e999f' }
				],
				colors: {
					'editor.foreground': '#333333',
					'editor.background': '#f8f8f8',
					'editorCursor.foreground': '#333333',
					'editor.lineHighlightBackground': '#eeeeee',
					'editor.selectionBackground': '#d6d6d6'
				}
			});

			editor = monaco.editor.create(editorEl, {
				value: value,
				language: 'python',
				theme: 'monokai-light',
				automaticLayout: true,
				bracketPairColorization: {
					enabled: true
				}
			});

			editor.onDidChangeModelContent(() => {
				const currentValue = editor!.getValue();
				if (currentValue !== value) {
					isUpdatingFromParent = true;
					value = currentValue;
					isUpdatingFromParent = false;
				}
			});
		}

		initializeMonaco();

		return () => {
			if (editor) {
				editor.dispose();
			}
		};
	});

	$effect(() => {
		value;

		const currentEditor = editor;
		if (currentEditor && value !== currentEditor.getValue() && !isUpdatingFromParent) {
			currentEditor.setValue(value);
		}
	});
</script>

<div class="monaco-editor-wrapper" bind:this={editorEl}></div>

<style>
	.monaco-editor-wrapper {
		width: 100%;
		height: 100%;
		min-height: 400px;
	}
</style>
<script lang="ts">
	import katex from 'katex';
	import 'katex/dist/katex.min.css'; // Re-enable the CSS
	import { parseLatex } from '../utils';
	import type { LatexSegment } from '../types';

	let { text = '' }: { text: string } = $props();

	let segments: LatexSegment[] = $derived(parseLatex(text));

	function renderMath(math: string): string {
		try {
			return katex.renderToString(math, {
				displayMode: false,
				throwOnError: true
			});
		} catch (error: any) {
			console.error('KaTeX rendering error:', error);
			return `<span class="katex-error">Error: ${error.message}</span>`;
		}
	}
</script>

<span class="wrapper">
	{#each segments as segment}
		{#if segment.type === 'text'}
			{@html segment.content}
		{:else if segment.type === 'latex'}
			<div class="katex-container">
				{@html renderMath(segment.content)}
			</div>
		{/if}
	{/each}
</span>

<style>
	.katex-error {
		color: red;
		font-family: monospace;
	}
    .wrapper {
        font-family: 'Lexend';
    }
	.katex-container {
		display: inline-block;
		position: relative;
	}
</style>


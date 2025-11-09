<script lang="ts">
	import katex from 'katex';
	import 'katex/dist/katex.min.css';
	import { parseLatex } from '$lib/utils';
	import type { LatexSegment } from '$lib/types';

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

<span>
	{#each segments as segment}
		{#if segment.type === 'text'}
			{segment.content}
		{:else if segment.type === 'latex'}
			{@html renderMath(segment.content)}
		{/if}
	{/each}
</span>

<style>
	.katex-error {
		color: red;
		font-family: monospace;
	}
    span {
        font-family: 'Lexend';
    }
</style>


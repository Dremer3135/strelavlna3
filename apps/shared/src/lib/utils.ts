import type { LatexSegment } from "./types.ts";

export function parseLatex(text: string): LatexSegment[] {
    const segments: LatexSegment[] = [];
    const regex = /\$(.*?)\$/g;
    let lastIndex = 0;
    let match;

    while ((match = regex.exec(text)) !== null) {
        // Push the text before the match
        if (match.index > lastIndex) {
            segments.push({ type: 'text', content: text.substring(lastIndex, match.index) });
        }
        // Push the latex part, using the captured group
        segments.push({ type: 'latex', content: match[1] });
        lastIndex = regex.lastIndex;
    }

    // Push any remaining text after the last match
    if (lastIndex < text.length) {
        segments.push({ type: 'text', content: text.substring(lastIndex) });
    }

    return segments;
}
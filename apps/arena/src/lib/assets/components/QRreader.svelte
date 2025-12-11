<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { readBarcodes } from "zxing-wasm/reader";
    import type { ReaderOptions } from "zxing-wasm/reader";

    // PROPS
    let { onScan } = $props<{ onScan: (text: string) => void }>();

    // STATE (Initialized to undefined to fix TS errors)
    let videoEl = $state<HTMLVideoElement>();
    let canvasEl = $state<HTMLCanvasElement>();
    
    let stream: MediaStream | null = null;
    let isScanning = false;
    let statusMessage = $state("Initializing camera..."); // To debug visibility

    // OPTIONS: Optimized for Data Matrix
    const scanOptions: ReaderOptions = {
        formats: ["DataMatrix", "QRCode"], 
        tryHarder: true, // Crucial for rotation/tilt
        maxNumberOfSymbols: 1
    };

    async function startCamera() {
        if (!videoEl) return;

        try {
            statusMessage = "Requesting permissions...";
            
            // 1. Request High Resolution & Focus
            stream = await navigator.mediaDevices.getUserMedia({ 
                video: { 
                    facingMode: "environment",
                    width: { ideal: 1920 }, // Request 1080p for better Data Matrix detail
                    height: { ideal: 1080 },
                    // @ts-ignore: standard constraint but TS sometimes complains
                    advanced: [{ focusMode: "continuous" }] 
                } 
            });

            statusMessage = "Starting video...";
            
            // 2. Attach stream to video element
            videoEl.srcObject = stream;
            
            // 3. Wait for video to be ready
            videoEl.onloadedmetadata = () => {
                videoEl?.play().then(() => {
                    statusMessage = ""; // Clear message on success
                    isScanning = true;
                    requestAnimationFrame(scanFrame);
                }).catch(e => {
                    console.error("Play error:", e);
                    statusMessage = "Error playing video.";
                });
            };

        } catch (err) {
            console.error("Camera Error:", err);
            statusMessage = "Camera access denied or not supported.";
        }
    }

    async function scanFrame() {
        if (!isScanning || !videoEl || !canvasEl) return;
        
        // Ensure video has actual data
        if (videoEl.readyState === videoEl.HAVE_ENOUGH_DATA) {
            const width = videoEl.videoWidth;
            const height = videoEl.videoHeight;

            // Sync canvas size
            canvasEl.width = width;
            canvasEl.height = height;

            const ctx = canvasEl.getContext('2d', { willReadFrequently: true });
            
            if (ctx) {
                // Draw current frame
                ctx.drawImage(videoEl, 0, 0, width, height);
                const imageData = ctx.getImageData(0, 0, width, height);

                try {
                    // Attempt Read
                    const results = await readBarcodes(imageData, scanOptions);
                    
                    if (results.length > 0) {
                        const text = results[0].text;
                        onScan(text);
                        // Optional: debounce or pause here if needed
                    }
                } catch (err) {
                    // WASM errors or no code found
                }
            }
        }

        if (isScanning) requestAnimationFrame(scanFrame);
    }

    onMount(() => {
        startCamera();
    });

    onDestroy(() => {
        isScanning = false;
        if (stream) {
            stream.getTracks().forEach(track => track.stop());
        }
    });
</script>

<div class="scanner-container">
    <!-- Feedback Message Overlay -->
    {#if statusMessage}
        <div class="status-overlay">{statusMessage}</div>
    {/if}

    <!-- 
      1. 'playsinline' is REQUIRED for iOS 
      2. 'muted' is REQUIRED for Chrome autoplay 
    -->
    <video 
        bind:this={videoEl} 
        playsinline 
        muted 
        autoplay
    ></video>
    
    <!-- Hidden Canvas for processing -->
    <canvas bind:this={canvasEl} style="display: none;"></canvas>
</div>

<style>
    .scanner-container {
        position: relative;
        width: 100%;
        height: 100%;
        /* IMPORTANT: Give it a minimum height or it will collapse to 0px */
        /* min-height: 400px; 
        background-color: #000; */
        /* overflow: hidden;
        border-radius: 8px; */
    }

    video {
        width: 100%;
        height: 100%;
        object-fit: cover; /* Ensures video fills the box without stretching */
        display: block;
    }

    .status-overlay {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: white;
        background: rgba(0,0,0,0.7);
        padding: 1rem;
        border-radius: 8px;
        z-index: 10;
        text-align: center;
    }
</style>
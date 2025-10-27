
<script>
    import Navbar from '$lib/components/general/navbar.svelte';
    import Footer from '$lib/components/general/footer.svelte';
    let { children } = $props();
</script>

<Navbar />
<main>
    {@render children?.()}
</main>
<Footer />


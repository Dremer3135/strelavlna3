
<script>
    import Navbar from '$lib/components/general/Navbar.svelte';
    import Footer from '$lib/components/general/footer.svelte';
    let { children } = $props();
</script>

<Navbar />
<main>
    {@render children?.()}
</main>
<Footer />

<style lang="scss">
    main {
        position: relative;
    }
</style>

<script lang="ts">
    import { onMount } from "svelte";
    
    let { data } = $props();




    onMount(() => {
        console.log(data.user);
    });
</script>


<main>
    <h1>Nothing here yet!</h1>
    <p>{data.user}</p>
</main>



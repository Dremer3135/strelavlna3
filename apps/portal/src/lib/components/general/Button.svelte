<script lang="ts">
  // We now also accept `children`, which is the Svelte 5 equivalent of a slot.
  let { class: className = '', children, ...rest } = $props();
</script>

<button
  class="base-button {className}"
  {...rest}
>
  <!-- We render the `children` snippet here. -->
  {@render children()}
</button>

<style>
  .base-button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 8px;
    background-color: #007bff;
    color: white;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease-in-out;
  }

  .base-button:hover {
    background-color: #0056b3;
  }

  .base-button:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }
</style>

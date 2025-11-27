<script lang="ts">
  // We now also accept `children`, which is the Svelte 5 equivalent of a slot.
  let { theme, disabled, children, ...rest }: { theme : "yellow" | "orange" | "purple" | "pink", disabled: boolean, children: any } = $props();
</script>

<button
  style="--color-theme: var(--color-{theme});"
  {...rest}
  class:disabled={disabled}
>
  <!-- We render the `children` snippet here. -->
  {@render children()}
</button>

<style lang="scss">
  button {
    all: unset;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 8px;
    border: 3px var(--color-theme) solid;
    border-radius: 3px;
    width: fit-content;
    cursor: pointer;
    padding: 5px 20px;
    background-color: color-mix(in srgb, var(--color-theme) 5%, transparent 95%);
    color: #333333;
    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

    &:hover {
      background-color: color-mix(in srgb, var(--color-theme) 20%, transparent 80%);
      color: color-mix(in srgb, var(--color-theme) 20%, black 80%);
    }

    &.disbaled {
      cursor: not-allowed;
      color: #555555;
      border-color: #AAAAAA;
      background-color: color-mix(in srgb, #333333 5%, transparent 95%);
    }
  }


</style>

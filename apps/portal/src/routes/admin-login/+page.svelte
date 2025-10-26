<script lang="ts">
    let { data, form } = $props();
</script>

<h1>Admin Login</h1>

<form method="POST">
    <input type="hidden" name="redirectTo" value={data.redirectTo} />

    <label>
        Email
        <input name="email" type="email" required>
    </label>
    <label>
        Password
        <input name="password" type="password" required>
    </label>

    {#if form?.error}
        <p style="color: red;">{form.error}</p>
    {/if}

    <button>Login</button>
</form>

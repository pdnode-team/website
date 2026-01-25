
<script>
  import CheckStatus from "$lib/components/CheckPocketbaseStatus.svelte"
  let { children } = $props();
</script>


<CheckStatus />
{@render children()}
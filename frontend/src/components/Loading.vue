<template>
    <div>
        <svg width="51px" height="50px" viewBox="0 0 51 50">
            <rect y="0" width="13" height="50" fill="#4527A0">
                <animate attributeName="height" values="50;10;50" begin="0s" dur="1s" repeatCount="indefinite" />
                <animate attributeName="y" values="0;20;0" begin="0s" dur="1s" repeatCount="indefinite" />
            </rect>
            <rect x="19" y="0" width="13" height="50" fill="#303F9F">
                <animate attributeName="height" values="50;10;50" begin="0.2s" dur="1s" repeatCount="indefinite" />
                <animate attributeName="y" values="0;20;0" begin="0.2s" dur="1s" repeatCount="indefinite" />
            </rect>
            <rect x="38" y="0" width="13" height="50" fill="#1E88E5">
                <animate attributeName="height" values="50;10;50" begin="0.4s" dur="1s" repeatCount="indefinite" />
                <animate attributeName="y" values="0;20;0" begin="0.4s" dur="1s" repeatCount="indefinite" />
            </rect>
        </svg>
    </div>
</template>

<script>
export default {
    name: "loader",
    data() {
        return {}
    },
};
</script>


<style lang="scss">
.styledSvg {
    width: 150px;
    height: 150px;
    margin: auto;
    display: block;
}
</style>
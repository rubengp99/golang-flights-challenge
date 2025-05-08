export default {
    beforeLeave(element) {
        this.prevHeight = getComputedStyle(element).height;
    },
    enter(element) {
        const { height } = getComputedStyle(element);

        element.style.height = this.prevHeight;

        setTimeout(() => {
            element.style.height = height;
        });
    },
    afterEnter(element) {
        element.style.height = 'auto';
    },
    animate(DEFAULT_TRANSITION) {
        this.$router.beforeEach((to, from, next) => {
            let transitionName = to.meta.transitionName || from.meta.transitionName;

            if (transitionName === 'slide') {
                const toDepth = to.path.split('/').length;
                const fromDepth = from.path.split('/').length;
                transitionName = toDepth < fromDepth ? 'slide-right' : 'slide-left';
            }

            this.transitionName = transitionName || DEFAULT_TRANSITION;

            next();
        });
    }
}
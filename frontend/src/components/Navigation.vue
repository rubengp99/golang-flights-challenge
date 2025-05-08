<template>
    <nav>
        <v-navigation-drawer v-model="drawer" fixed temporary>
            <v-list-item two-line>
                <v-list-item-avatar class="no-radius">
                    <v-img :src="require('@/assets/logo.svg')"></v-img>
                </v-list-item-avatar>

                <v-list-item-content>
                    <v-list-item-title class="title">Flight Finder</v-list-item-title>
                    <v-list-item-subtitle>Find the best pricings for your vacations</v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>

            <v-divider></v-divider>

            <v-list dense nav>
                <v-list-item v-for="item in items" :key="item.title" :to="item.path">
                    <v-list-item-avatar class="no-radius">
                        <v-img width="32px" :src="item.icon"></v-img>
                    </v-list-item-avatar>

                    <v-list-item-content>
                        <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-app-bar class="h-59" :clipped-left="'mx-5'" app color="darken-4" dark>
            <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
            <v-spacer></v-spacer>
            <v-toolbar-title style="width: 300px" class="ml-0 pl-4">
                <span>
                    <v-btn icon large>
                        <v-avatar size="32px" item>
                            <v-img :src="require('@/assets/logo2.svg')"></v-img>
                        </v-avatar>
                    </v-btn>Flight Finder
                </span>
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <div v-if="!user.loggedIn">
                <v-btn text to="/login" class="mx-5 font-weight-bold white--text text-capitalize">Log In</v-btn>
            </div>
        </v-app-bar>
    </nav>
</template>

<script>
import { mapState } from "vuex";

export default {
    props: {
        source: String
    },
    data: () => ({
        dialog: false,
        drawer: false,
        items: [
            { icon: require("@/assets/home.svg"), title: "Dashboard", path: "/" },
        ]
    }),
    computed: {
        ...mapState(["user"])
    }
};
</script>

<style lang="scss">
.no-radius {
    border-radius: 0;
}

.h-59 {
    min-height: 59px !important;
}
</style>
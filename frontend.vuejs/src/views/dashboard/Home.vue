<template>
    <div class="home">
        <v-container class="mt-15">
            <v-form v-model="valid">
                <v-container>
                    <v-row cols="15">
                        <v-col cols="12" md="3">
                            <v-select prepend-inner-icon="mdi-map-marker" v-model="origin" label="Origin"
                                :items="locations" item-title="label" item-value="value" />
                        </v-col>

                        <v-col cols="12" md="3">
                            <v-select prepend-inner-icon="mdi-map-marker" v-model="destination" label="Destiny"
                                :items="filteredDestinations" item-title="label" item-value="value" />
                        </v-col>

                        <v-col cols="12" md="2">
                            <v-menu v-model="menu" :close-on-content-click="false" transition="scale-transition"
                                offset-y min-width="290px">
                                <template #activator="{ props }">
                                    <v-text-field v-bind="props" v-model="formattedDate" label="Date" readonly
                                        prepend-inner-icon="mdi-calendar" />
                                </template>
                                <v-date-picker v-model="date" :min="minDate" @update:model-value="onDateSelect"
                                    color="primary" />
                            </v-menu>
                        </v-col>
                        <v-col cols="12" md="2">
                            <v-text-field v-model="quantity" label="Adults" type="number" min="1" step="1"
                                prepend-inner-icon="mdi-account" hide-details />
                        </v-col>
                        <v-col cols="12" md="2">
                            <v-btn class="mt-2" color="primary" rounded="0" @click="searchFlight"
                                :disabled="!canSearch">
                                <v-icon start>mdi-magnify</v-icon>
                                Search
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-container>
            </v-form>
            <v-container>
                <v-row>
                    <v-col cols="12" md="6">
                        <h3 class="text-h6 mb-2">Cheapest Flights</h3>
                        <v-data-table :headers="headers" :items="cheapestFlights" :items-per-page="15" item-value="id"
                            :loading="isLoading" class="elevation-1">
                            <template #item="{ item }">
                                <tr @click="showDetails(item)" style="cursor: pointer">
                                    <td>{{ item.airline }}</td>
                                    <td>{{ item.flightNumber }}</td>
                                    <td>${{ item.price.value.toFixed(2) }} {{ item.price.currency }}</td>
                                    <td>{{ item.durationInMinutes }} min</td>
                                </tr>
                            </template>
                        </v-data-table>
                    </v-col>

                    <v-col cols="12" md="6">
                        <h3 class="text-h6 mb-2">Fastest Flights</h3>
                        <v-data-table :headers="headers" :items="fastestFlights" :items-per-page="15" item-value="id"
                            :loading="isLoading" class="elevation-1">
                            <template #item="{ item }">
                                <tr @click="showDetails(item)" style="cursor: pointer">
                                    <td>{{ item.airline }}</td>
                                    <td>{{ item.flightNumber }}</td>
                                    <td>${{ item.price.value.toFixed(2) }} {{ item.price.currency }}</td>
                                    <td>{{ item.durationInMinutes }} min</td>
                                </tr>
                            </template>
                        </v-data-table>
                    </v-col>
                </v-row>

                <v-dialog v-model="dialog" max-width="600">
                    <v-card>
                        <v-card-title class="text-h6">
                            Flight Details: {{ selectedFlight?.flightNumber }}
                        </v-card-title>
                        <v-card-text v-if="selectedFlight && selectedFlight.flightNumber">
                            <p><strong>Airline:</strong> {{ selectedFlight.airline }}</p>
                            <p><strong>Flight Number:</strong> {{ selectedFlight.flightNumber }}</p>
                            <p><strong>Departure:</strong> {{ selectedFlight.departure.timestamp }} from {{
                                selectedFlight.departure.iataCode }}</p>
                            <p><strong>Arrival:</strong> {{ selectedFlight.arrival.timestamp }} at {{
                                selectedFlight.arrival.iataCode }}</p>
                            <p><strong>Duration:</strong> {{ selectedFlight.durationInMinutes }} minutes</p>
                            <p><strong>Layovers:</strong> {{ selectedFlight.layovers }}</p>
                            <p><strong>Price:</strong> {{ selectedFlight.price.value }} {{ selectedFlight.price.currency
                                }}</p>
                        </v-card-text>
                        <v-card-actions>
                            <v-spacer />
                            <v-btn color="primary" @click="dialog = false">Close</v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </v-container>
        </v-container>
    </div>
</template>

<script>
import { toast } from 'vue3-toastify';
import backend from '@/services/backend';

const tomorrow = new Date();
tomorrow.setDate(tomorrow.getDate() + 1);
const minDate = tomorrow.toISOString().split('T')[0]; // "YYYY-MM-DD"
export default {
    name: "Home",
    data() {
        return {
            minDate: minDate,
            menu: false,
            date: null,
            valid: false,
            quantity: 1,
            origin: null,
            destination: null,
            selectedSort: 'Price',
            flights: {
                cheapest: [],
                fastest: [],
            },
            dialog: false,
            isLoading: false,
            selectedFlight: null,
            headers: [
                { title: 'Airline', key: 'airline', sortable: false },
                { title: 'Flight Number', key: 'flightNumber', sortable: false },
                {
                    title: 'Price',
                    key: 'price.value',
                    value: (item) => `$${item.price.value.toFixed(2)} ${item.price.currency}`,
                    sortable: false
                },
                {
                    title: 'Duration',
                    key: 'durationInMinutes',
                    value: (item) => `${item.durationInMinutes} min`,
                    sortable: false
                },
            ],
            locations: [
                { label: 'Los Angeles (LAX)', value: 'LAX' },
                { label: 'New York (JFK)', value: 'JFK' },
                { label: 'London Heathrow (LHR)', value: 'LHR' },
                { label: 'Paris Charles de Gaulle (CDG)', value: 'CDG' },
                { label: 'Frankfurt (FRA)', value: 'FRA' },
                { label: 'Tokyo Haneda (HND)', value: 'HND' },
                { label: 'Singapore Changi (SIN)', value: 'SIN' },
                { label: 'Sydney (SYD)', value: 'SYD' },
                { label: 'Dubai (DXB)', value: 'DXB' },
                { label: 'Toronto Pearson (YYZ)', value: 'YYZ' },
                { label: 'São Paulo (GRU)', value: 'GRU' },
                { label: 'Hong Kong (HKG)', value: 'HKG' },
                { label: 'Seoul Incheon (ICN)', value: 'ICN' },
                { label: 'Bangkok (BKK)', value: 'BKK' },
                { label: 'Istanbul (IST)', value: 'IST' }
            ]
        };
    },
    computed: {
        formattedDate() {
            return this.date
                ? new Date(this.date).toLocaleDateString()
                : '';
        },
        filteredOrigins() {
            return this.locations.filter(loc => loc.value !== this.destination);
        },
        filteredDestinations() {
            return this.locations.filter(loc => loc.value !== this.origin);
        },
        cheapestFlights() {
            return Array.isArray(this.flights.cheapest)
                ? this.flights.cheapest.map(f => ({
                    ...f,
                    id: `${f.airline}-${f.flightNumber}`
                }))
                : [];
        },
        fastestFlights() {
            return Array.isArray(this.flights.fastest)
                ? this.flights.fastest.map(f => ({
                    ...f,
                    id: `${f.airline}-${f.flightNumber}`
                }))
                : [];
        },
        canSearch() {
            return (
                this.origin &&
                this.destination &&
                this.date &&
                Number(this.quantity) >= 1
            )
        }
    },
    methods: {
        searchFlight() {
            const formatDate = (date) => {
                const d = new Date(date);
                const year = d.getFullYear();
                const month = String(d.getMonth() + 1).padStart(2, '0');
                const day = String(d.getDate()).padStart(2, '0');
                return `${year}-${month}-${day}`;
            };

            this.isLoading = true
            const params = {
                origin: this.origin,
                destination: this.destination,
                date: formatDate(this.date),  // ensure ISO format
                adults: this.quantity,
            };
            const token = sessionStorage.getItem("token");

            backend().get("flights/search",
                {
                    params,
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }).then((response) => {
                    this.isLoading = false,
                        this.flights = response.data
                    toast.success("Nice, we have some flights for you!");
                }).catch(e => {
                    console.log(e);
                    toast.error("Something went wrong!");
                });
        },
        onDateSelect() {
            this.menu = false;
        },
        showDetails(item) {
            this.selectedFlight = item;
            this.dialog = true;
        },
    },
};
</script>

<style lang="scss">
.v-data-table-header th {
    cursor: default !important;
}

.absolute-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}
</style>
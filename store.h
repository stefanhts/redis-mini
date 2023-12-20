struct Value {
    char* str;
};

struct Key {
    char* key;
};

struct Pair {
    struct char* key;
    struct char* value;
};

struct Store {
    struct Pair* items;
    size_t len;
    size_t capacity;
};

struct Store* initStore(size_t len) {
    struct Store *s = malloc(sizeof(struct Pair) * len);
    return s;
}

void freeStore(struct Store *s) {
    for(int i = 0; i < s->len; i++) {
        free(&s->items[i]);
    }
    free(s);
}

void insert(struct Store* s, char* key, char* value) {
    struct *Pair p = (*Pair)malloc(sizeof(Pair));
    p->key = (*char)malloc(strlen(key));
    p->value = (*char)malloc(strlen(value));
    strcpy(p->key, key);
    strcpy(p->value, value);
    // Check if the length of the store with a new element will exceed the allotted lenght, if so, allocate more space
    if (s->len >= s->capacity -1) {
       growStore(&s);
    }

    // Put the item at end of list
    s->items[s->len] = p;
    // Increment len
    s->len = s->len + 1;
}

void growStore(struct Store **s) {
    // TODO: implement this
    printf("NOT IMPLEMENTED YET: `growStore()`");
}

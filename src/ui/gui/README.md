# GUI

ターミナルで動かすためのUI。

[Meditator Pattern](https://refactoring.guru/ja/design-patterns/mediator) を採用して、Meditator ↔︎ Component のやり取りのみを許可する。Component ↔︎ Component のやり取りを禁止することで、Component 同士が複雑に絡み合い、依存し合うスパゲッティコードを避ける。

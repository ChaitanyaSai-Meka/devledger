CREATE TABLE User (
        UserID INTEGER PRIMARY KEY,
        UserName TEXT   
);

CREATE TABLE Groups (
        GroupID INTEGER PRIMARY KEY,
        GroupName TEXT  
);

CREATE TABLE GroupMember (
        GroupID INTEGER,   
        UserID INTEGER,   
        PRIMARY KEY (GroupID , UserID),
        FOREIGN KEY (GroupID) REFERENCES Groups(GroupID),
        FOREIGN KEY (UserID) REFERENCES User(UserID)
);

CREATE TABLE Expense (
	ExpenseID INTEGER PRIMARY KEY,
	Amount REAL,
	Description TEXT,
	PaidByUserID INTEGER,
	GroupID INTEGER,
	FOREIGN KEY (PaidByUserID) REFERENCES User(UserID),
	FOREIGN KEY (GroupID) REFERENCES Groups(GroupID)
);

CREATE TABLE Split (
    ExpenseID INTEGER,
	UserID INTEGER,
	Amount REAL,
    PRIMARY KEY (ExpenseID, UserID),
	Settled BOOLEAN,
    FOREIGN KEY (ExpenseID) REFERENCES Expense(ExpenseID),
    FOREIGN KEY (UserID) REFERENCES User(UserID)
);


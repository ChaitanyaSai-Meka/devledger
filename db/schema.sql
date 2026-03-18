CREATE TABLE Users (
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
    FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

CREATE TABLE Expense (
	ExpenseID INTEGER PRIMARY KEY,
	Amount INTEGER,
	Description TEXT,
	PaidByUserID INTEGER,
	GroupID INTEGER,
	FOREIGN KEY (PaidByUserID) REFERENCES Users(UserID),
	FOREIGN KEY (GroupID) REFERENCES Groups(GroupID)
);

CREATE TABLE Split (
    ExpenseID INTEGER,
	UserID INTEGER,
	Amount INTEGER,
    PRIMARY KEY (ExpenseID, UserID),
	Settled BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (ExpenseID) REFERENCES Expense(ExpenseID),
    FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

